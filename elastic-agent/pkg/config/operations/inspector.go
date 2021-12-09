// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package operations

import (
	"fmt"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/application/info"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/application/paths"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/configuration"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/errors"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/storage"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/storage/store"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/config"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/core/logger"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/fleetapi"
)

var (
	// ErrNoFleetConfig is returned when no configuration was retrieved from fleet just yet.
	ErrNoFleetConfig = fmt.Errorf("no fleet config retrieved yet")
)

// LoadFullAgentConfig load agent config based on provided paths and defined capabilities.
// In case fleet is used, config from policy action is returned.
func LoadFullAgentConfig(cfgPath string, failOnFleetMissing bool) (*config.Config, error) {
	rawConfig, err := loadConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg, err := configuration.NewFromConfig(rawConfig)
	if err != nil {
		return nil, err
	}

	if configuration.IsStandalone(cfg.Fleet) {
		return rawConfig, nil
	}

	fleetConfig, err := loadFleetConfig(rawConfig)
	if err != nil {
		return nil, err
	} else if fleetConfig == nil {
		if failOnFleetMissing {
			return nil, ErrNoFleetConfig
		}

		// resolving fleet config but not fleet config retrieved yet, returning last applied config
		return rawConfig, nil
	}

	return config.NewConfigFrom(fleetConfig)
}

func loadConfig(configPath string) (*config.Config, error) {
	rawConfig, err := config.LoadFile(configPath)
	if err != nil {
		return nil, err
	}

	path := paths.AgentConfigFile()

	store := storage.NewDiskStore(path)
	reader, err := store.Load()
	if err != nil {
		return nil, errors.New(err, "could not initialize config store",
			errors.TypeFilesystem,
			errors.M(errors.MetaKeyPath, path))
	}

	config, err := config.NewConfigFrom(reader)
	if err != nil {
		return nil, errors.New(err,
			fmt.Sprintf("fail to read configuration %s for the elastic-agent", path),
			errors.TypeFilesystem,
			errors.M(errors.MetaKeyPath, path))
	}

	// merge local configuration and configuration persisted from fleet.
	rawConfig.Merge(config)

	if err := info.InjectAgentConfig(rawConfig); err != nil {
		return nil, err
	}

	return rawConfig, nil
}

func loadFleetConfig(cfg *config.Config) (map[string]interface{}, error) {
	log, err := newErrorLogger()
	if err != nil {
		return nil, err
	}

	stateStore, err := store.NewStateStoreWithMigration(log, paths.AgentActionStoreFile(), paths.AgentStateStoreFile())
	if err != nil {
		return nil, err
	}

	for _, c := range stateStore.Actions() {
		cfgChange, ok := c.(*fleetapi.ActionPolicyChange)
		if !ok {
			continue
		}

		return cfgChange.Policy, nil
	}
	return nil, nil
}

func newErrorLogger() (*logger.Logger, error) {
	return logger.NewWithLogpLevel("", logp.ErrorLevel, false)
}
