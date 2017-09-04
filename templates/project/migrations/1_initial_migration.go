package migrations

import (
	"{{.project}}/bindings"
	. "github.com/swarmdotmarket/perigord/migration"
)

var initial_migration = &Migration{
	Number: 1,
	F: func(m *Migrator) error {
		_, _, _, err := bindings.DeployMigrations(m.Auth(), m.Backend())
		if err != nil {
			return err
		}

		m.Backend().Commit()

		return nil
	},
}

func init() {
	AddMigration(initial_migration)
}
