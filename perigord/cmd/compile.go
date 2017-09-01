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
	"strings"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile contract source files",
	Run: func(cmd *cobra.Command, args []string) {
		err := runInRoot(func() error {
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

func init() {
	RootCmd.AddCommand(compileCmd)
}

func compileContracts() error {
	// TODO: Figure out relative imports and if we need to do anything else here
	matches, err := filepath.Glob(ContractsDirectory + "/*.sol")
	if err != nil {
		return err
	}

	for _, match := range matches {
		compileContract(match)
	}

	return nil
}

func compileContract(path string) error {
	// TODO: This just shells out atm, could directly integrate abigen and call
	// into it as a library later
	command := "solc"
	args := []string{path, "--bin", "--abi", "--optimize", "--overwrite", "-o", BuildDirectory}
	return execWithOutput(command, args...)
}

func generateBindings() error {
	matches, err := filepath.Glob(BuildDirectory + "/*.abi")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(BindingsDirectory, os.FileMode(0755)); err != nil {
		return err
	}

	for _, match := range matches {
		generateBinding(strings.TrimSuffix(match, filepath.Ext(match)))
	}

	return nil
}

func generateBinding(path string) error {
	// TODO: Allow alternate binding directories / package names, in config file
	command := "abigen"
	name := filepath.Base(path)
	abifile := path + ".abi"
	binfile := path + ".bin"
	outfile := filepath.Join(BindingsDirectory, filepath.Base(name)) + ".go"
	args := []string{"--abi", abifile, "--bin", binfile, "--pkg", "bindings", "--type", name, "--out", outfile}
	return execWithOutput(command, args...)
}
