// Copyright © 2017 Swarm Market <info@swarm.market>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
