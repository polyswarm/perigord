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
