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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ProjectConfigFilename = "perigord.yaml"
	ContractsDirectory    = "contracts"
	BuildDirectory        = "build"
	BindingsDirectory     = "bindings"
	MigrationsDirectory   = "migrations"
	TestsDirectory        = "tests"
)

func Fatal(v ...interface{}) {
	fmt.Println("Error:", v)
	os.Exit(1)
}

// This function from github.com/spf13/cobra used under Apache 2.0 license
func exists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if !os.IsNotExist(err) {
		return false, err
	}

	return false, nil
}

// This function from github.com/spf13/cobra used under Apache 2.0 license
func isEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		Fatal(err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		Fatal(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		Fatal(err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}

// Public as we also run from the stub
func FindProject() (*Project, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return findProject(wd)
}

func findProject(path string) (*Project, error) {
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return nil, errors.New("Could not find project root")
	}

	e, err := exists(filepath.Join(path, ProjectConfigFilename))
	if err != nil {
		return nil, err
	}
	if e {
		return NewProjectFromPath(path), nil
	}

	return findProject(filepath.Dir(path))
}

func ExecWithOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExecWithPipes(command string, in []byte, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(in)
	}()

	return cmd.CombinedOutput()
}

func RunInRoot(f func() error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	project, err := findProject(wd)
	if err != nil {
		return err
	}

	if err := os.Chdir(project.AbsPath()); err != nil {
		return err
	}
	defer os.Chdir(wd)

	return f()
}

func runStub(stubcommand string, stubargs ...string) error {
	return RunInRoot(func() error {
		// Our stub depends on up to date bindings, make sure they're generated
		command := runtime.GOROOT() + "/bin/go"
		args := []string{"generate"}
		err := ExecWithOutput(command, args...)
		if err != nil {
			return err
		}

		args = append([]string{"run", "stub/main.go", stubcommand}, stubargs...)
		return ExecWithOutput(command, args...)
	})
}
