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

package migration

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

type MigrationFunc func(context.Context, *bind.TransactOpts, bind.ContractBackend) error

type Migration struct {
	Number int
	F      MigrationFunc
}

type Migrations []*Migration

func (s Migrations) Len() int {
	return len(s)
}

func (s Migrations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Migrations) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

type Migrator struct {
	auth       *bind.TransactOpts
	backend    bind.ContractBackend
	migrations Migrations
}

var migrator *Migrator = &Migrator{}

func (m *Migrator) Auth() *bind.TransactOpts {
	return m.auth
}

func (m *Migrator) Backend() bind.ContractBackend {
	return m.backend
}

func (m *Migrator) AddMigration(migration *Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) RunMigrations(ctx context.Context) error {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(10000000000)}})

	m.auth = auth
	m.backend = backend

	// TODO: Check migration contract for last run and only run new
	sort.Sort(m.migrations)
	for _, migration := range m.migrations {
		fmt.Println("Running migration", migration.Number)
		if err := migration.F(ctx, m.auth, m.backend); err != nil {
			return err
		}
	}

	backend.Commit()

	return nil
}

func Auth() *bind.TransactOpts {
	return migrator.Auth()
}

func Backend() bind.ContractBackend {
	return migrator.Backend()
}

func AddMigration(migration *Migration) {
	migrator.AddMigration(migration)
}

func RunMigrations(ctx context.Context) error {
	return migrator.RunMigrations(ctx)
}
