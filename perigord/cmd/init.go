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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/swarmdotmarket/perigord/templates"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new Ethereum project with example contracts and tests",
	Long: `Initialize (perigord init) will create a new application, with a license
and the appropriate structure for a perigord-based project.

  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/swarmdotmarket/perigord);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}

		var project *Project
		if len(args) == 0 {
			project = NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				project = NewProjectFromPath(arg)
			} else {
				project = NewProject(arg)
			}
		} else {
			Fatal("please provide only one argument")
		}

		initializeProject(project)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	initCmd.PersistentFlags().StringP("license", "l", "", "name of license for the project")

	viper.BindPFlag("author", initCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.BindPFlag("license", initCmd.PersistentFlags().Lookup("license"))
	viper.SetDefault("license", "apache")
}

func initializeProject(project *Project) {
	if err := os.MkdirAll(project.AbsPath(), os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	if !isEmpty(project.AbsPath()) {
		Fatal("Cowardly refusing to initialize a non-empty dir")
	}

	if err := templates.RestoreTemplates(project.AbsPath(), "project", "project", project.TemplateData()); err != nil {
		Fatal(err)
	}
}
