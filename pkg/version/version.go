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

package version

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Application build information.
var (
	Branch    string
	BuildDate string
	GitSHA1   string
	Version   = "v0.1.0"
)

// Print writes application version details to standard output.
func Print() {
	fmt.Printf("hierarchy, version %v (branch: %v, revision: %v), build date: %v, go version: %v\n", Version, Branch, GitSHA1, BuildDate, runtime.Version())
}

// Print writes application version details to the log
func Log() {
	log.WithFields(log.Fields{
		"version":    Version,
		"branch":     Branch,
		"revision":   GitSHA1,
		"build date": BuildDate,
		"go version": runtime.Version(),
	}).Info("Hierarchy")
}