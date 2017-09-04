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
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/swarmdotmarket/perigord/templates"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new contract to the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			Fatal("Must specify package name")
		}

		name := args[0]

		// TODO: Allow full contract paths or detecting filenames
		match, _ := regexp.MatchString("\\w+", name)
		if !match {
			Fatal("Invalid contract name specified")
		}

		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}

		project, err := FindProject(wd)
		if err != nil {
			Fatal(err)
		}

		addContract(name, project)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

func addContract(name string, project *Project) {
	path := filepath.Join(project.AbsPath(), ContractsDirectory, name+".sol")

	data := map[string]string{"contract": name}

	if err := templates.RestoreTemplate(path, "contract/contract.sol", data); err != nil {
		Fatal(err)
	}
}
