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

package fleetapi

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActionSerialization(t *testing.T) {
	a := ActionApp{
		ActionID:   "1231232",
		ActionType: "APP_INPUT",
		InputType:  "osquery",
		Data:       []byte(`{ "foo": "bar" }`),
	}

	m, err := a.MarshalMap()
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(4, len(m))
	if diff != "" {
		t.Error(diff)
	}

	diff = cmp.Diff(a.ActionID, mapStringVal(m, "id"))
	if diff != "" {
		t.Error(diff)
	}

	diff = cmp.Diff(a.ActionType, mapStringVal(m, "type"))
	if diff != "" {
		t.Error(diff)
	}

	diff = cmp.Diff(a.InputType, mapStringVal(m, "input_type"))
	if diff != "" {
		t.Error(diff)
	}

	diff = cmp.Diff(a.Data, mapRawMessageVal(m, "data"))
	if diff != "" {
		t.Error(diff)
	}

	diff = cmp.Diff(a.StartedAt, mapStringVal(m, "started_at"))
	if diff != "" {
		t.Error(diff)
	}
	diff = cmp.Diff(a.CompletedAt, mapStringVal(m, "completed_at"))
	if diff != "" {
		t.Error(diff)
	}
	diff = cmp.Diff(a.Error, mapStringVal(m, "error"))
	if diff != "" {
		t.Error(diff)
	}
}

func mapStringVal(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func mapRawMessageVal(m map[string]interface{}, key string) json.RawMessage {
	if v, ok := m[key]; ok {
		if res, ok := v.(json.RawMessage); ok {
			return res
		}
	}
	return nil
}
