// Copyright © 2017 PolySwarm <info@polyswarm.io>
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
	"fmt"
	"os"
	"path/filepath"

	"github.com/polyswarm/perigord/project"
	"github.com/polyswarm/perigord/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new Ethereum project with example contracts and tests",
	Long: `Initialize (perigord init) will create a new application, with a license
and the appropriate structure for a perigord-based project.

  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/polyswarm/perigord);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}

		var prj *project.Project
		if len(args) == 0 {
			prj = project.NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				prj = project.NewProjectFromPath(arg)
			} else {
				prj = project.NewProject(arg)
			}
		} else {
			Fatal("please provide only one argument")
		}

		initializeProject(prj)
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

func initializeProject(prj *project.Project) {
	if err := os.MkdirAll(prj.AbsPath(), os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	if !isEmpty(prj.AbsPath()) {
		Fatal("Cowardly refusing to initialize a non-empty dir")
	}

	if err := templates.RestoreTemplates(prj.AbsPath(), "project", "project", prj.TemplateData()); err != nil {
		Fatal(err)
	}

	fmt.Println("Project initialized in", prj.AbsPath())
}
