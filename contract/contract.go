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
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractDeployer interface {
	Deploy(context.Context, *bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, interface{}, error)
	Bind(context.Context, *bind.TransactOpts, bind.ContractBackend, common.Address) (interface{}, error)
}

type Contract struct {
	Address  common.Address
	Session  interface{}
	deployed bool
	deployer ContractDeployer
}

func (c *Contract) Deploy(ctx context.Context, auth *bind.TransactOpts, backend bind.ContractBackend) error {
	// TODO: Is this the correct behavior?
	if !c.deployed {
		address, transaction, session, err := c.deployer.Deploy(ctx, auth, backend)
		if err != nil {
			return err
		}

		deployBackend, ok := backend.(bind.DeployBackend)
		if ok {
			address, err = bind.WaitDeployed(ctx, deployBackend, transaction)
			if err != nil {
				return err
			}
		}

		c.Address = address
		c.Session = session
		c.deployed = true
		return nil
	} else {
		session, err := c.deployer.Bind(ctx, auth, backend, c.Address)
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

func Deploy(ctx context.Context, name string, auth *bind.TransactOpts, backend bind.ContractBackend) error {
	contract := contracts[name]
	if contract == nil {
		return errors.New("No such contract found")
	}

	return contract.Deploy(ctx, auth, backend)
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
