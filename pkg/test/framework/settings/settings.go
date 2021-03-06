//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package settings

import (
	"fmt"
	"strings"

	"istio.io/istio/pkg/test/env"

	"github.com/google/uuid"
)

// EnvironmentID is a unique identifier for a testing environment.
type EnvironmentID string

const (
	// MaxTestIDLength is the maximum length allowed for testID.
	MaxTestIDLength = 30

	// Local environment identifier
	Local = EnvironmentID("local")

	// Kubernetes environment identifier
	Kubernetes = EnvironmentID("kubernetes")
)

var (
	// ISTIO_TEST_ENVIRONMENT indicates in which mode the test framework should run. It can be "local", or
	// "kube".
	// nolint: golint
	ISTIO_TEST_ENVIRONMENT env.Variable = "ISTIO_TEST_ENVIRONMENT"

	globalSettings = &Settings{
		Environment: EnvironmentID(ISTIO_TEST_ENVIRONMENT.ValueOrDefault(string(Local))),
	}
)

// Settings is the set of arguments to the test driver.
type Settings struct {
	// Environment to run the tests in. By default, a local environment will be used.
	Environment EnvironmentID

	// TestID is the id of the test suite. This should supplied by the author once, and must be immutable.
	TestID string

	// RunID is the id of the current run.
	RunID string

	// Do not cleanup the resources after the test run.
	NoCleanup bool

	// Local working directory root for creating temporary directories / files in. If left empty,
	// os.TempDir() will be used.
	WorkDir string
}

// New returns settings built from flags and environment variables.
func New(testID string) (*Settings, error) {
	// Copy the default settings.
	s := &(*globalSettings)

	s.TestID = testID
	s.RunID = generateRunID(testID)

	if err := s.validate(); err != nil {
		return nil, err
	}

	return s, nil
}

// Validate the arguments.
func (a *Settings) validate() error {
	switch a.Environment {
	case Local, Kubernetes:
	default:
		return fmt.Errorf("unrecognized environment: %q", string(a.Environment))
	}

	if a.TestID == "" || len(a.TestID) > MaxTestIDLength {
		return fmt.Errorf("testID must be non-empty and cannot be longer than %d characters", MaxTestIDLength)
	}

	return nil
}

func generateRunID(testID string) string {
	u := uuid.New().String()
	u = strings.Replace(u, "-", "", -1)
	testID = strings.Replace(testID, "_", "-", -1)
	// We want at least 6 characters of uuid padding
	padding := MaxTestIDLength - len(testID)
	return fmt.Sprintf("%s-%s", testID, u[0:padding])
}
