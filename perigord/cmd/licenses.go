// Copyright © 2015 Steve Francia <spf@spf13.com>.
// Modifications Copyright © 2017 PolySwarm <info@polyswarm.io>
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

// Parts inspired by https://github.com/ryanuber/go-license

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/polyswarm/perigord/templates"
)

// Licenses contains all possible licenses a user can choose from.
var Licenses = make(map[string]License)

// License represents a software license agreement, containing the Name of
// the license, its possible matches (on the command line as given to cobra),
// the header to be used with each file on the file's creating, and the text
// of the license
type License struct {
	Name            string   // The type of license in use
	PossibleMatches []string // Similar names to guess
	Text            string   // License text data
	Header          string   // License header for source files
}

func init() {
	// Allows a user to not use a license.
	Licenses["none"] = License{"None", []string{"none", "false"}, "", ""}

	initializeLicense("agpl", "GNU Affero General Public License", []string{"agpl", "affero gpl", "gnu agpl"})
	initializeLicense("apache", "Apache 2.0", []string{"apache", "apache20", "apache 2.0", "apache2.0", "apache-2.0"})
	initializeLicense("freebsd", "Simplified BSD License", []string{"freebsd", "simpbsd", "simple bsd", "2-clause bsd", "2 clause bsd", "simplified bsd license"})
	initializeLicense("bsd", "NewBSD", []string{"bsd", "newbsd", "3 clause bsd", "3-clause bsd"})
	initializeLicense("gpl2", "GNU General Public License 2.0", []string{"gpl2", "gnu gpl2", "gplv2"})
	initializeLicense("gpl3", "GNU General Public License 3.0", []string{"gpl3", "gplv3", "gpl", "gnu gpl3", "gnu gpl"})
	initializeLicense("lgpl", "GNU Lesser General Public License", []string{"lgpl", "lesser gpl", "gnu lgpl"})
	initializeLicense("mit", "MIT License", []string{"mit"})
}

func initializeLicense(asset, name string, possibleMatches []string) error {
	data := map[string]string{
		"copyright": copyrightLine(),
	}

	buf, err := templates.ExecuteTemplate(fmt.Sprintf("licenses/%s/header.tpl", asset), data)
	if err != nil {
		return err
	}
	header := buf.String()

	buf, err = templates.ExecuteTemplate(fmt.Sprintf("licenses/%s/text.tpl", asset), data)
	if err != nil {
		return err
	}
	text := buf.String()

	Licenses[asset] = License{name, possibleMatches, header, text}
	return nil
}

// getLicense returns license specified by user in flag or in config.
// If user didn't specify the license, it returns Apache License 2.0.
//
// TODO: Inspect project for existing license
func getLicense() License {
	// If user wants to have custom license, use that.
	if viper.IsSet("license.header") || viper.IsSet("license.text") {
		return License{Header: viper.GetString("license.header"),
			Text: "license.text"}
	}

	// If user wants to have built-in license, use that.
	if viper.IsSet("license") {
		return findLicense(viper.GetString("license"))
	}

	// If user didn't set any license, use Apache 2.0 by default.
	return Licenses["apache"]
}

func copyrightLine() string {
	author := viper.GetString("author")

	year := viper.GetString("year") // For tests.
	if year == "" {
		year = time.Now().Format("2006")
	}

	return "Copyright © " + year + " " + author
}

// findLicense looks for License object of built-in licenses.
// If it didn't find license, then the app will be terminated and
// error will be printed.
func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		Fatal("unknown license: " + name)
	}
	return Licenses[found]
}

// matchLicense compares the given a license name
// to PossibleMatches of all built-in licenses.
// It returns blank string, if name is blank string or it didn't find
// then appropriate match to name.
func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
