// Copyright © 2017 Swarm Market <info@swarm.market>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	perigord "github.com/swarmdotmarket/perigord/perigord/cmd"
)

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		err := perigord.RunInRoot(func() error {
			command := runtime.GOROOT() + "/bin/go"
			args := []string{"test"}
			return perigord.ExecWithOutput(command, args...)
		})
		if err != nil {
			perigord.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
