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

package errors

// ErrorType defines an error
type ErrorType int

const (
	// TypeUnexpected is a default error type for errors without specified type.
	TypeUnexpected ErrorType = iota
	// TypeConfig is a configuration error.
	TypeConfig
	// TypePath is an invalid path error.
	TypePath
	// TypeApplication is an error describing errors related to application state.
	TypeApplication
	// TypeApplicationCrash is an error describing unexpected application crash.
	TypeApplicationCrash
	// TypeNetwork represents collection of errors related to networking.
	TypeNetwork
	// TypeFilesystem represents set of errors generated by filesystem operations.
	TypeFilesystem
	// TypeSecurity represents set of errors related to security, encryption, etc.
	TypeSecurity
)

const (
	// MetaKeyPath is a metadata key used used in filesystem errors.
	MetaKeyPath = "path"
	// MetaKeyURI is a metadata key used in network related errors.
	MetaKeyURI = "uri"
	// MetaKeyAppID is a metadata key used to identify application related to error.
	MetaKeyAppID = "app_id"
	// MetaKeyAppName is a metadata key used to specify application name related to error.
	MetaKeyAppName = "app_name"
)

var readableTypes = map[ErrorType]string{
	TypeUnexpected:       "UNEXPECTED",
	TypeConfig:           "CONFIG",
	TypePath:             "PATH",
	TypeApplicationCrash: "CRASH",
	TypeApplication:      "APPLICATION",
	TypeNetwork:          "NETWORK",
	TypeFilesystem:       "FILESYSTEM",
	TypeSecurity:         "SECURITY",
}
