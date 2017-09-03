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

package migration

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

type MigrationFunc func(*Migrator) error

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
	Auth       *bind.TransactOpts
	Backend    *backends.SimulatedBackend
	migrations Migrations
}

func (m *Migrator) RunMigrations() error {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	m.Auth = bind.NewKeyedTransactor(key)
	m.Backend = backends.NewSimulatedBackend(core.GenesisAlloc{m.Auth.From: {Balance: big.NewInt(10000000000)}})

	// TODO: Check migration contract for last run and only run new
	sort.Sort(m.migrations)
	for _, migration := range m.migrations {
		if err := migration.F(m); err != nil {
			return err
		}
	}

	m.Backend.Commit()

	return nil
}

var migrator *Migrator = &Migrator{}

func (m *Migrator) AddMigration(migration *Migration) {
	m.migrations = append(m.migrations, migration)
}

func AddMigration(migration *Migration) {
	migrator.AddMigration(migration)
}

func RunMigrations() error {
	return migrator.RunMigrations()
}
