// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/palantir/amalgomate/amalgomated"
	amalgomatedcheck "github.com/palantir/godel-okgo-asset-golint/generated_src"
	"github.com/palantir/godel-okgo-asset-golint/golint/config"
	"github.com/palantir/godel-okgo-asset-golint/golint/creator"
	"github.com/palantir/okgo/checker"
	"github.com/palantir/pkg/cobracli"
)

const (
	cpuProfilePublicFlagName  = "cpuprofile"
	cpuProfilePrivateFlagName = "cpuprofile-private"
)

func main() {
	// if cpuProfilePrivateFlagName flag is present, profile
	if cpuProfileFlagVal, osArgs := getFlagVal(cpuProfilePrivateFlagName, os.Args); cpuProfileFlagVal != "" {
		f, err := os.Create(cpuProfileFlagVal)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to create CPU profile: %v\n", err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
		os.Args = osArgs
	}

	os.Exit(amalgomated.RunApp(os.Args, nil, amalgomated.NewCmdLibrary(amalgomatedcheck.Instance()), checkMain))
}

func checkMain(osArgs []string) int {
	var cpuProfileFlagVal string
	cpuProfileFlagVal, osArgs = getFlagVal(cpuProfilePublicFlagName, osArgs)
	os.Args = osArgs
	rootCmd := checker.AssetRootCmd(creator.Golint(cpuProfileFlagVal), config.UpgradeConfig, "run golint check")
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func getFlagVal(flagName string, osArgs []string) (flagVal string, osArgsWithFlagAndValRemoved []string) {
	for idx, arg := range osArgs {
		if arg == "--"+flagName && idx+1 < len(osArgs) {
			flagVal = osArgs[idx+1]
			osArgs = append(osArgs[:idx], osArgs[idx+2:]...)
			break
		}
	}
	return flagVal, osArgs
}
