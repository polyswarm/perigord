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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new Ethereum project with example contracts and tests",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			Fatal("Must specify package name")
		}

		name := args[0]

		// TODO: Allow full package paths or init-ing a directory like cobra
		match, _ := regexp.MatchString("\\w+", name)
		if !match {
			Fatal("Invalid package name specified")
		}

		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}

		path := filepath.Join(wd, name)

		initProject(name, path)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// TODO: Add options like defining package names, etc
	// Can add to config yaml and parse from project root in later invocations
}

func initProject(name, path string) {
	if err := os.MkdirAll(path, os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	data := map[string]string{"project": name}
	if err := templates.ExecuteTemplates(path, "project", "project", data); err != nil {
		Fatal(err)
	}
}
