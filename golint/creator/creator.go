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

package creator

import (
	"os"
	"os/exec"
	"strings"

	"github.com/palantir/godel-okgo-asset-golint/golint"
	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"
)

func Golint() checker.Creator {
	return checker.NewCreator(
		golint.TypeName,
		golint.Priority,
		func(cfgYML []byte) (okgo.Checker, error) {
			// set GOROOT environment variable if it is not set
			if _, ok := os.LookupEnv("GOROOT"); !ok {
				if gorootOutput, err := exec.Command("go", "env", "GOROOT").CombinedOutput(); err == nil {
					_ = os.Setenv("GOROOT", strings.TrimSpace(string(gorootOutput)))
				}
			}
			return checker.NewAmalgomatedChecker(golint.TypeName, checker.ParamPriority(golint.Priority)), nil
		},
	)
}
