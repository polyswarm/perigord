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

package project

import (
	"errors"
	"os"
	"path/filepath"
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
