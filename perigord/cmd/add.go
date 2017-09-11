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
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/swarmdotmarket/perigord/templates"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new contract or test to the project",
}

var addContractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Add a new contract to the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			Fatal("Must specify contract name")
		}

		name := args[0]

		match, _ := regexp.MatchString("\\w+", name)
		if !match {
			Fatal("Invalid contract name specified")
		}

		project, err := FindProject()
		if err != nil {
			Fatal(err)
		}

		addContract(name, project)
	},
}

var addMigrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Add a new migration to the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			Fatal("Must specify migration name")
		}

		name := args[0]

		match, _ := regexp.MatchString("\\w+", name)
		if !match {
			Fatal("Invalid test name specified")
		}

		project, err := FindProject()
		if err != nil {
			Fatal(err)
		}

		addMigration(name, project)
	},
}

var addTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Add a new test to the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			Fatal("Must specify test name")
		}

		name := args[0]

		match, _ := regexp.MatchString("\\w+", name)
		if !match {
			Fatal("Invalid test name specified")
		}

		project, err := FindProject()
		if err != nil {
			Fatal(err)
		}

		addTest(name, project)
	},
}

func init() {
	addCmd.AddCommand(addContractCmd)
	addCmd.AddCommand(addMigrationCmd)
	addCmd.AddCommand(addTestCmd)
	RootCmd.AddCommand(addCmd)
}

func addContract(name string, project *Project) {
	path := filepath.Join(project.AbsPath(), ContractsDirectory, name+".sol")

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	data := project.TemplateData()
	data["contract"] = name

	if err := templates.RestoreTemplate(path, "contract/contract.sol", data); err != nil {
		Fatal(err)
	}

	fmt.Println("New contract added at", path)
}

func addMigration(name string, project *Project) {
	path := filepath.Join(project.AbsPath(), MigrationsDirectory)
	glob, err := filepath.Glob(filepath.Join(path, "*.go"))

	numMigrations := 1
	if err == nil {
		numMigrations += len(glob)
	}

	path = filepath.Join(path, fmt.Sprintf("%d_%s.go", numMigrations, name))

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	data := project.TemplateData()
	data["contract"] = name
	data["number"] = numMigrations

	if err := templates.RestoreTemplate(path, "migration/migration.go", data); err != nil {
		Fatal(err)
	}

	fmt.Println("New migration added at", path)
}

func addTest(name string, project *Project) {
	path := filepath.Join(project.AbsPath(), TestsDirectory, name+".go")

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		Fatal(err)
	}

	data := project.TemplateData()
	data["test"] = name

	if err := templates.RestoreTemplate(path, "test/test.go", data); err != nil {
		Fatal(err)
	}

	fmt.Println("New test added at", path)
}
