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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	fmt.Println("Error: ", v)
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
		command := "go"
		args := append([]string{"run", "stub/main.go", stubcommand}, stubargs...)
		return execWithOutput(command, args...)
	})
}
