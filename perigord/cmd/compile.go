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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/polyswarm/perigord/templates"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile contract source files",
	Run: func(cmd *cobra.Command, args []string) {
		err := RunInRoot(func() error {
			if err := compileContracts(); err != nil {
				return err
			}

			return generateBindings()
		})
		if err != nil {
			Fatal(err)
		}
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "(alias for compile)",
	Run:   compileCmd.Run,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "(alias for compile)",
	Run:   compileCmd.Run,
}

func init() {
	RootCmd.AddCommand(compileCmd)
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(generateCmd)
}

func compileContracts() error {
	if err := os.Chdir(ContractsDirectory); err != nil {
		return err
	}

	matches := make([]string, 0)
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".sol") {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(matches) == 0 {
		Fatal("No contracts found, create one with `perigord add contract NAME`")
	}

	command := "solc"
	_, err = exec.LookPath(command)
	if err != nil {
		return errors.New("Can't locate solc, is it installed and in your path")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	compilerConfig, err := templates.ExecuteTemplate("solc/solc.json.tpl", matches)
	if err != nil {
		return err
	}

	args := []string{"--allow-paths", dir, "--standard-json"}
	outputJson, err := ExecWithPipes(command, compilerConfig.Bytes(), args...)
	if err != nil {
		return err
	}

	if err = os.Chdir(".."); err != nil {
		return err
	}

	var output map[string]interface{}
	if err = json.Unmarshal(outputJson, &output); err != nil {
		return err
	}

	if err = handleErrors(output); err != nil {
		return err
	}

	return saveArtifacts(output)
}

func handleErrors(output map[string]interface{}) error {
	hasFatal := false
	compilerErrors, ok := output["errors"].([]interface{})
	if !ok {
		return errors.New("Invalid json")
	}

	for _, value := range compilerErrors {
		compilerErr, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("Invalid json")
		}

		severity, ok := compilerErr["severity"].(string)
		if !ok {
			return errors.New("Invalid json")
		}

		message, ok := compilerErr["formattedMessage"].(string)
		if !ok {
			return errors.New("Invalid json")
		}

		if severity != "warning" {
			hasFatal = true
		}

		fmt.Println(message)
	}

	if hasFatal {
		return errors.New("Error detected, aborting. Please check solidity output for more details.")
	}

	return nil
}

func saveArtifacts(output map[string]interface{}) error {
	if err := os.MkdirAll(BuildDirectory, os.FileMode(0755)); err != nil {
		return err
	}

	contracts, ok := output["contracts"].(map[string]interface{})
	if !ok {
		return errors.New("Invalid json")
	}

	for _, value := range contracts {
		contract, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("Invalid json")
		}

		for name, value := range contract {
			data, ok := value.(map[string]interface{})
			if !ok {
				return errors.New("Invalid json")
			}

			// Get ABI as a json blob
			abi, err := json.Marshal(data["abi"])
			if err != nil {
				return err
			}

			evm, ok := data["evm"].(map[string]interface{})
			if !ok {
				return errors.New("Invalid json")
			}

			bytecodeObject, ok := evm["bytecode"].(map[string]interface{})
			if !ok {
				return errors.New("Invalid json")
			}

			// Get bytecode as a bin
			bytecode, ok := bytecodeObject["object"].(string)
			if !ok {
				return errors.New("Invalid json")
			}

			// Get link references as a json blob
			linkReferences, err := json.Marshal(bytecodeObject["linkReferences"])
			if err != nil {
				return err
			}

			ioutil.WriteFile(filepath.Join(BuildDirectory, name+".abi"), abi, 0644)
			ioutil.WriteFile(filepath.Join(BuildDirectory, name+".bin"), []byte(bytecode), 0644)
			ioutil.WriteFile(filepath.Join(BuildDirectory, name+".link"), linkReferences, 0644)
		}
	}

	return nil
}

func generateBindings() error {
	return nil
	//	matches, err := filepath.Glob(BuildDirectory + "/*.abi")
	//	if err != nil {
	//		return err
	//	}
	//
	//	if err := os.MkdirAll(BindingsDirectory, os.FileMode(0755)); err != nil {
	//		return err
	//	}
	//
	//	for _, match := range matches {
	//		if err := generateBinding(strings.TrimSuffix(match, filepath.Ext(match))); err != nil {
	//			return err
	//		}
	//	}
	//
	//	return nil
}

func generateBinding(path string) error {
	return nil
	//	// TODO: Allow alternate binding directories / package names, in config file
	//	command := "abigen"
	//	_, err := exec.LookPath(command)
	//	if err != nil {
	//		return errors.New("Can't locate abigen, is it installed and in your path")
	//	}
	//
	//	name := filepath.Base(path)
	//	abifile := path + ".abi"
	//	binfile := path + ".bin"
	//	outfile := filepath.Join(BindingsDirectory, filepath.Base(name)) + ".go"
	//	args := []string{"--abi", abifile, "--bin", binfile, "--pkg", "bindings", "--type", name, "--out", outfile}
	//	return ExecWithOutput(command, args...)
}
