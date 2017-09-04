package migrations

// TODO: Remove relative import and templatize package name
import (
	"fmt"
	"math/big"

	"{{.project}}/bindings"
	. "github.com/swarmdotmarket/perigord/migration"
)

var migration = &Migration{
	Number: 1,
	F: func(m *Migrator) error {
		_, _, interactor, err := bindings.DeployMigrations(m.Auth, m.Backend)
		if err != nil {
			return err
		}

		m.Backend.Commit()

		t, _ := interactor.SetCompleted(m.Auth, new(big.Int))
		fmt.Println("SetCompleted:", t)

		return nil
	},
}

func init() {
	AddMigration(migration)
}
