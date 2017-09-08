// Example main file for a native dapp, replace with application code
package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/swarmdotmarket/perigord/migration"

	"{{.project}}/bindings"
)

func main() {
	// Run our migrations
	migration.RunMigrations()

	// Deploy Foo contract manually
	// TODO: This will be handled by contract management code within perigord
	// lib in an upcoming revision
	_, _, interactor, _ := bindings.DeployFoo(migration.Auth(), migration.Backend())
	migration.Backend().Commit()

	ret, _ := interactor.Bar(&bind.CallOpts{Pending: true})
	fmt.Println("Foo returned", ret.Int64())
}
