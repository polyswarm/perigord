// Example main file for a native dapp, replace with application code
package main

import (
	"context"

	"github.com/polyswarm/perigord/migration"

	_ "{{.project}}/migrations"
)

func main() {
	// Run our migrations
	migration.RunMigrations(context.Background())
}
