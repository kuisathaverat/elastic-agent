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

package configuration

import (
	"net/url"

	"github.com/elastic/beats/v7/libbeat/common/transport/tlscommon"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/errors"
)

// FleetServerConfig is the configuration written so Elastic Agent can run Fleet Server.
type FleetServerConfig struct {
	Bootstrap    bool                     `config:"bootstrap" yaml:"bootstrap,omitempty"`
	Policy       *FleetServerPolicyConfig `config:"policy" yaml:"policy,omitempty"`
	Output       FleetServerOutputConfig  `config:"output" yaml:"output,omitempty"`
	Host         string                   `config:"host" yaml:"host,omitempty"`
	Port         uint16                   `config:"port" yaml:"port,omitempty"`
	InternalPort uint16                   `config:"internal_port" yaml:"internal_port,omitempty"`
	TLS          *tlscommon.Config        `config:"ssl" yaml:"ssl,omitempty"`
}

// FleetServerPolicyConfig is the configuration for the policy Fleet Server should run on.
type FleetServerPolicyConfig struct {
	ID string `config:"id"`
}

// FleetServerOutputConfig is the connection for Fleet Server to call to Elasticsearch.
type FleetServerOutputConfig struct {
	Elasticsearch Elasticsearch `config:"elasticsearch" yaml:"elasticsearch"`
}

// Elasticsearch is the configuration for elasticsearch.
type Elasticsearch struct {
	Protocol     string            `config:"protocol" yaml:"protocol"`
	Hosts        []string          `config:"hosts" yaml:"hosts"`
	Path         string            `config:"path" yaml:"path,omitempty"`
	ServiceToken string            `config:"service_token" yaml:"service_token,omitempty"`
	TLS          *tlscommon.Config `config:"ssl" yaml:"ssl,omitempty"`
	Headers      map[string]string `config:"headers" yaml:"headers,omitempty"`
	ProxyURL     string            `config:"proxy_url" yaml:"proxy_url,omitempty"`
	ProxyDisable bool              `config:"proxy_disable" yaml:"proxy_disable"`
	ProxyHeaders map[string]string `config:"proxy_headers" yaml:"proxy_headers"`
}

// ElasticsearchFromConnStr returns an Elasticsearch configuration from the connection string.
func ElasticsearchFromConnStr(conn string, serviceToken string, insecure bool) (Elasticsearch, error) {
	u, err := url.Parse(conn)
	if err != nil {
		return Elasticsearch{}, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return Elasticsearch{}, errors.New("invalid connection string: scheme must be http or https")
	}
	if u.Host == "" {
		return Elasticsearch{}, errors.New("invalid connection string: must include a host")
	}
	cfg := Elasticsearch{
		Protocol: u.Scheme,
		Hosts:    []string{u.Host},
		Path:     u.Path,
		TLS:      nil,
	}
	if insecure {
		cfg.TLS = &tlscommon.Config{
			VerificationMode: tlscommon.VerifyNone,
		}
	}
	if serviceToken == "" {
		return Elasticsearch{}, errors.New("invalid connection string: must include a service token")
	}
	cfg.ServiceToken = serviceToken
	return cfg, nil
}
