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
	"sort"

	"github.com/polyswarm/perigord/network"
)

type MigrationFunc func(context.Context, *network.Network) error

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
	migrations Migrations
}

var migrator *Migrator = &Migrator{}

func (m *Migrator) AddMigration(migration *Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) RunMigrations(ctx context.Context, net *network.Network) error {
	sort.Sort(m.migrations)
	for _, migration := range m.migrations {
		fmt.Println("Running migration", migration.Number)
		if err := migration.F(ctx, net); err != nil {
			return err
		}
	}

	return nil
}

func AddMigration(migration *Migration) {
	migrator.AddMigration(migration)
}

func RunMigrations(ctx context.Context, net *network.Network) error {
	return migrator.RunMigrations(ctx, net)
}
