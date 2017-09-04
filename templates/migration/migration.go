package migrations

import (
	"{{.project}}/bindings"
	. "github.com/swarmdotmarket/perigord/migration"
)

var {{.migration}} = &Migration{
	Number: {{.number}},
	F: func(m *Migrator) error {
		return nil
	},
}

func init() {
	AddMigration({{.migration}})
}
