// Copyright {yyyy} {name of copyright owner}
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
//
// This file contains a modified version of kubectl's HandlePluginCommand function:
// [https://github.com/kubernetes/kubectl/blob/d0e936c703558421df2ef367a7c8c1c5972e5c3c/pkg/cmd/cmd.go#L260-L307]
//
// As well as DefaultPluginHandler.Lookup method:
// [https://github.com/kubernetes/kubectl/blob/d0e936c703558421df2ef367a7c8c1c5972e5c3c/pkg/cmd/cmd.go#L209-L219]

package kubectl

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	validPluginPrefixes = []string{"kubectl"}
)

type PluginHandler interface {
	Lookup(filename string) (string, bool)
}

// IsPlugin receives a pluginHandler and command-line arguments and attempts to find
// a plugin executable on the PATH that satisfies the given arguments.
func IsPlugin(cmdArgs []string, pluginHandler PluginHandler) bool {
	var remainingArgs []string // all "non-flag" arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.Replace(arg, "-", "_", -1))
	}

	if len(remainingArgs) == 0 {
		// the length of cmdArgs is at least 1
		return false // flags cannot be placed before plugin name
	}

	// attempt to find binary, starting at longest possible name with given cmdArgs
	for len(remainingArgs) > 0 {
		_, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		// found it! It's a plugin alright
		return true
	}

	// No plugins found
	return false
}

type DefaultPluginHandler struct{}

// Ensure it implements the interface
var _ PluginHandler = DefaultPluginHandler{}

// Lookup implements PluginHandler
func (DefaultPluginHandler) Lookup(filename string) (string, bool) {
	for _, prefix := range validPluginPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if shouldSkipOnLookPathErr(err) || len(path) == 0 {
			continue
		}
		return path, true
	}
	return "", false
}

func shouldSkipOnLookPathErr(err error) bool {
	return err != nil && !errors.Is(err, exec.ErrDot)
}
