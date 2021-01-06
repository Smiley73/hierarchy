/*
Copyright 2021 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetFilesSuccess verifies that we receive the correct list of files to be merged
// As part of this test, it will also check the proper functioning of the regex for the file filter
// fail.txt and fail.yaml.disabled should never be returned
func TestGetFilesSuccess(t *testing.T) {
	expected := []string{"testdata/default/defaults.json", "testdata/default/defaults.yml"}
	result := getFiles("testdata/default", defaultFileFilter)
	assert.Equal(t, expected, result)

	expected = []string{"testdata/yaml/one.yaml", "testdata/yaml/two.yml"}
	result = getFiles("testdata/yaml", defaultFileFilter)
	assert.Equal(t, expected, result)
}

// TestProcessHierarchySuccess verifies that all directories are correctly added to the list to be processed
// It will also test the correct handling of comments and different ways of specifying a relative path
func TestProcessHierarchySuccess(t *testing.T) {
	var cfg config

	cfg.hierarchyFile = "testdata/test1/hierarchy.lst"
	cfg.basePath = "testdata/test1"
	cfg.outputFile = "output.yaml"
	cfg.filterExtension = defaultFileFilter
	cfg.logDebug = false
	cfg.logTrace = false
	cfg.failMissing = false

	expected := []string{"testdata/default", "testdata/yaml", "testdata/json", "testdata/empty", "testdata/test1"}
	result := processHierarchy(cfg)
	assert.Equal(t, expected, result)
}

// TestFailMissing tests the correct behavior of the `--failmissing` command line option
// It spawns a new process to determine the exit code of the application.
// Anything other than a 1 is a problem
// It uses the environment variable TEST_FAIL_EMPTY to signal the actual execution of the functionality
func TestFailMissing(t *testing.T) {
	if os.Getenv("TEST_FAIL_EMPTY") == "1" {
		var cfg config
		cfg.hierarchyFile = "testdata/test1/hierarchy.lst"
		cfg.basePath = "testdata/test1"
		cfg.outputFile = "output.yaml"
		cfg.filterExtension = defaultFileFilter
		cfg.logDebug = false
		cfg.logTrace = false
		cfg.failMissing = true

		processHierarchy(cfg)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFailMissing")
	cmd.Env = append(os.Environ(), "TEST_FAIL_EMPTY=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		fmt.Printf("Process correctly failed with %v\n", e)
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestEnd2EndSuccess runs through the full functionality end-to-end
// It compares the generated final file with one stored in git
func TestEnd2EndSuccess(t *testing.T) {
	var cfg config

	cfg.hierarchyFile = "testdata/test1/hierarchy.lst"
	cfg.basePath = "testdata/test1"
	cfg.outputFile = "output.yaml"
	cfg.filterExtension = defaultFileFilter
	cfg.logDebug = false
	cfg.logTrace = false
	cfg.failMissing = false

	// process the hierarchy and get the list of include files
	hierarchy := processHierarchy(cfg)

	// Lets do the deed
	mergeFilesInHierarchy(hierarchy, cfg.filterExtension, cfg.outputFile)

	expected, err := ioutil.ReadFile("testdata/test1/result/expected.yaml")
	if err != nil {
		t.Fatalf("Error reading file with expected test results: %v", err)
	}
	result, err := ioutil.ReadFile(cfg.outputFile)
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
	assert.Equal(t, string(expected), string(result))
}

// TestFailMissingEnvironmentVariable ensures that the application is correctly failing
// If an environment variable specified in `hierarchy.lst` is not found.
// It spawns a new process to determine the exit code of the application.
// Anything other than a 1 is a problem
// It uses the environment variable TEST_FAIL_EMPTY to signal the actual execution of the functionality
func TestFailMissingEnvironmentVariable(t *testing.T) {
	if os.Getenv("TEST_FAIL_EMPTY") == "1" {
		var cfg config
		cfg.hierarchyFile = "testdata/test2-with-env-fail/hierarchy.lst"
		cfg.basePath = "testdata/test2-with-env-fail"
		cfg.outputFile = "output.yaml"
		cfg.filterExtension = defaultFileFilter
		cfg.logDebug = false
		cfg.logTrace = false
		cfg.failMissing = false

		processHierarchy(cfg)

		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFailMissingEnvironmentVariable")
	cmd.Env = append(os.Environ(), "TEST_FAIL_EMPTY=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		fmt.Printf("Process correctly failed with %v\n", e)
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestEnd2EndEnvironmentVariablesSuccess runs through the full functionality end-to-end,
// specifically testing the correct resolution of environment variables specified in `hierarchy.lst`
// It compares the generated final file with one stored in git
func TestEnd2EndEnvironmentVariablesSuccess(t *testing.T) {
	var cfg config

	cfg.hierarchyFile = "testdata/test2-with-env/hierarchy.lst"
	cfg.basePath = "testdata/test2-with-env"
	cfg.outputFile = "output.yaml"
	cfg.filterExtension = defaultFileFilter
	cfg.logDebug = false
	cfg.logTrace = false
	cfg.failMissing = false

	// set the test environment variable
	os.Setenv("JSON", "json")

	// process the hierarchy and get the list of include files
	hierarchy := processHierarchy(cfg)

	// Merge files
	mergeFilesInHierarchy(hierarchy, cfg.filterExtension, cfg.outputFile)

	expected, err := ioutil.ReadFile("testdata/test2-with-env/result/expected.yaml")
	if err != nil {
		t.Fatalf("Error reading file with expected test results: %v", err)
	}
	result, err := ioutil.ReadFile(cfg.outputFile)
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
	assert.Equal(t, string(expected), string(result))
}