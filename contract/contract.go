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

package contract

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/polyswarm/perigord/migration"
	"github.com/polyswarm/perigord/project"
)

type ContractDeployer interface {
	Deploy(context.Context, *migration.Network) (common.Address, *types.Transaction, interface{}, error)
	Bind(context.Context, *migration.Network, common.Address) (interface{}, error)
}

type Contract struct {
	Address  common.Address
	deployed bool
	Session  interface{}      `json:"-"`
	deployer ContractDeployer `json:"-"`
}

func (c *Contract) Deploy(ctx context.Context, network *migration.Network) error {
	if !c.deployed {
		address, _, session, err := c.deployer.Deploy(ctx, network)
		if err != nil {
			return err
		}

		backend := network.Backend()
		code, err := backend.CodeAt(ctx, address, nil)
		for err != nil || len(code) == 0 {
			time.Sleep(time.Second)
			code, err = backend.CodeAt(ctx, address, nil)
		}

		c.Address = address
		c.Session = session
		c.deployed = true
		return nil
	} else {
		session, err := c.deployer.Bind(ctx, network, c.Address)
		if err != nil {
			return err
		}

		c.Session = session
		return nil
	}
}

var contracts map[string]*Contract = make(map[string]*Contract)

func AddContract(name string, deployer ContractDeployer) {
	contracts[name] = &Contract{
		deployer: deployer,
	}
}

func Deploy(ctx context.Context, name string, network *migration.Network) error {
	contract := contracts[name]
	if contract == nil {
		return errors.New("No such contract found")
	}

	if err := contract.Deploy(ctx, network); err != nil {
		return err
	}

	if err := RecordDeployments(network); err != nil {
		return err
	}

	return nil
}

func RecordDeployments(network *migration.Network) error {
	project, err := project.FindProject()
	if err != nil {
		return err
	}

	data, err := json.Marshal(contracts)
	if err != nil {
		return err
	}

	network_path := filepath.Join(project.AbsPath(), network.Name()+".json")
	return ioutil.WriteFile(network_path, data, 0644)
}

func LoadDeployments(network *migration.Network) error {
	project, err := project.FindProject()
	if err != nil {
		return err
	}

	network_path := filepath.Join(project.AbsPath(), network.Name()+".json")
	data, err := ioutil.ReadFile(network_path)
	if err != nil {
		return err
	}

	var loaded_contracts map[string]*Contract
	if err := json.Unmarshal(data, &loaded_contracts); err != nil {
		return err
	}

	// Retain our initialized deployers, bind our sessions
	for name, contract := range loaded_contracts {
		contract.deployer = contracts[name].deployer
		contract.Deploy(context.Background(), network)
	}

	contracts = loaded_contracts
	return nil
}

func Session(name string) interface{} {
	contract := contracts[name]
	if contract == nil || !contract.deployed {
		return nil
	}
	return contract.Session
}

func Reset() {
	for k, v := range contracts {
		contracts[k] = &Contract{
			deployer: v.deployer,
		}
	}
}

func AddressOf(name string) common.Address {
	contract := contracts[name]
	if contract == nil || !contract.deployed {
		contract.Address{}
	}

	return contract.Address
}
