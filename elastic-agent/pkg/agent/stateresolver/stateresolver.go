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

package stateresolver

import (
	"sync"

	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/agent/configrequest"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/core/logger"
	uid "github.com/elastic/elastic-agent-poc/elastic-agent/pkg/id"
)

// Acker allow to ack the should state from a converge operation.
type Acker func()

// StateResolver is a resolver of a config state change
// it subscribes to Config event and publishes StateChange events based on that/
// Based on StateChange event operator know what to do.
type StateResolver struct {
	l        *logger.Logger
	curState state
	mu       sync.Mutex
}

// NewStateResolver allow to modify default event names.
func NewStateResolver(log *logger.Logger) (*StateResolver, error) {
	return &StateResolver{
		l: log,
	}, nil
}

// Resolve resolves passed config into one or multiple steps
func (s *StateResolver) Resolve(
	cfg configrequest.Request,
) (uid.ID, string, []configrequest.Step, Acker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newState, steps := converge(s.curState, cfg)
	newStateID := newState.ShortID()
	id, err := uid.Generate()
	if err != nil {
		return id, newStateID, nil, nil, err
	}

	s.l.Infof("New State ID is %s", newStateID)
	s.l.Infof("Converging state requires execution of %d step(s)", len(steps))
	for i, step := range steps {
		// more detailed debug log
		s.l.Debugf("step %d: %s", i, step.String())
	}

	// Allow the operator to ack the should state when applying the steps is done correctly.
	ack := func() {
		s.ack(newState)
	}

	return id, newStateID, steps, ack, nil
}

func (s *StateResolver) ack(newState state) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.l.Info("Updating internal state")
	s.curState = newState
}
