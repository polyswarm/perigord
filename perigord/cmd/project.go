// Copyright © 2015 Steve Francia <spf@spf13.com>.
// Modifications Copyright © 2017 Swarm Market <info@swarm.market>
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
	"runtime"
	"strings"
)

var srcPaths []string

func init() {
	// Initialize srcPaths.
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 {
		Fatal("$GOPATH is not set")
	}
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths {
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	}
}

// Project contains name, license and paths to projects.
type Project struct {
	absPath string
	srcPath string
	license License
	name    string
}

// NewProject returns Project with specified project name.
// If projectName is blank string, it returns nil.
func NewProject(projectName string) *Project {
	if projectName == "" {
		return nil
	}

	p := new(Project)
	p.name = projectName

	// 1. Find already created protect.
	p.absPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.absPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			Fatal(err)
		}
		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				p.absPath = filepath.Join(srcPath, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.absPath == "" {
		p.absPath = filepath.Join(srcPaths[0], projectName)
	}

	return p
}

// findPackage returns full path to existing go package in GOPATHs.
// findPackage returns "", if it can't find path.
// If packageName is "", findPackage returns "".
func findPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if e, _ := exists(packagePath); e {
			return packagePath
		}
	}

	return ""
}

// NewProjectFromPath returns Project with specified absolute path to
// package.
// If absPath is blank string or if absPath is not actually absolute,
// it returns nil.
func NewProjectFromPath(absPath string) *Project {
	if absPath == "" || !filepath.IsAbs(absPath) {
		return nil
	}

	p := new(Project)
	p.absPath = absPath
	p.name = filepath.ToSlash(trimSrcPath(p.absPath, p.SrcPath()))
	return p
}

// trimSrcPath trims at the beginning of absPath the srcPath.
func trimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		Fatal("Cobra supports project only within $GOPATH: " + err.Error())
	}
	return relPath
}

// License returns the License object of project.
func (p *Project) License() License {
	if p.license.Text == "" && p.license.Name != "None" {
		p.license = getLicense()
	}

	return p.license
}

// Name returns the name of project, e.g. "github.com/spf13/cobra"
func (p Project) Name() string {
	return p.name
}

// AbsPath returns absolute path of project.
func (p Project) AbsPath() string {
	return p.absPath
}

// SrcPath returns absolute path to $GOPATH/src where project is located.
func (p *Project) SrcPath() string {
	if p.srcPath != "" {
		return p.srcPath
	}

	if p.absPath == "" {
		p.srcPath = srcPaths[0]
		return p.srcPath
	}

	for _, srcPath := range srcPaths {
		if filepathHasPrefix(p.absPath, srcPath) {
			p.srcPath = srcPath
			break
		}
	}

	return p.srcPath
}

func (p *Project) TemplateData() map[string]interface{} {
	return map[string]interface{}{
		"project":   p.Name(),
		"license":   p.License(),
		"copyright": copyrightLine(),
	}
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}

	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}

	return path[0:len(prefix)] == prefix
}
