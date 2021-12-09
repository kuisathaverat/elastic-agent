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

package kubernetes

import (
	"time"

	"github.com/elastic/beats/v7/libbeat/common/kubernetes"
	"github.com/elastic/beats/v7/libbeat/common/kubernetes/metadata"
	"github.com/elastic/beats/v7/libbeat/logp"
)

// Config for kubernetes provider
type Config struct {
	Scope     string    `config:"scope"`
	Resources Resources `config:"resources"`

	KubeConfig        string                       `config:"kube_config"`
	KubeClientOptions kubernetes.KubeClientOptions `config:"kube_client_options"`

	Namespace      string        `config:"namespace"`
	SyncPeriod     time.Duration `config:"sync_period"`
	CleanupTimeout time.Duration `config:"cleanup_timeout" validate:"positive"`

	// Needed when resource is a Pod or Node
	Node string `config:"node"`

	AddResourceMetadata *metadata.AddResourceMetadataConfig `config:"add_resource_metadata"`
	IncludeLabels       []string                            `config:"include_labels"`
	ExcludeLabels       []string                            `config:"exclude_labels"`
	IncludeAnnotations  []string                            `config:"include_annotations"`

	LabelsDedot      bool `config:"labels.dedot"`
	AnnotationsDedot bool `config:"annotations.dedot"`
}

// Resources config section for resources' config blocks
type Resources struct {
	Pod     Enabled `config:"pod"`
	Node    Enabled `config:"node"`
	Service Enabled `config:"service"`
}

// Enabled config section for resources' config blocks
type Enabled struct {
	Enabled bool `config:"enabled"`
}

// InitDefaults initializes the default values for the config.
func (c *Config) InitDefaults() {
	c.CleanupTimeout = 60 * time.Second
	c.SyncPeriod = 10 * time.Minute
	c.Scope = "node"
	c.LabelsDedot = true
	c.AnnotationsDedot = true
	c.AddResourceMetadata = metadata.GetDefaultResourceMetadataConfig()
}

// Validate ensures correctness of config
func (c *Config) Validate() error {
	// Check if resource is service. If yes then default the scope to "cluster".
	if c.Resources.Service.Enabled {
		if c.Scope == "node" {
			logp.L().Warnf("can not set scope to `node` when using resource `Service`. resetting scope to `cluster`")
		}
		c.Scope = "cluster"
	}

	if !c.Resources.Pod.Enabled && !c.Resources.Node.Enabled && !c.Resources.Service.Enabled {
		c.Resources.Pod = Enabled{true}
		c.Resources.Node = Enabled{true}
	}

	return nil
}
