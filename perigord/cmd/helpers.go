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
	"errors"
	"fmt"
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
	TestDirectory         = "test"
)

func Fatal(v ...interface{}) {
	fmt.Println("Error:", v)
	os.Exit(1)
}

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

func FindRoot(path string) (string, error) {
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return "", errors.New("Could not find project root")
	}

	e, err := exists(filepath.Join(path, ProjectConfigFilename))
	if err != nil {
		return "", err
	}
	if e {
		return path, nil
	}

	return FindRoot(filepath.Dir(path))
}

func execWithOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runInRoot(f func() error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	root, err := FindRoot(wd)
	if err != nil {
		return err
	}

	if err := os.Chdir(root); err != nil {
		return err
	}
	defer os.Chdir(wd)

	return f()
}

func runStub(stubcommand string, stubargs ...string) error {
	return runInRoot(func() error {
		command := runtime.GOROOT() + "/bin/go"
		args := append([]string{"run", "stub/main.go", stubcommand}, stubargs...)
		return execWithOutput(command, args...)
	})
}
