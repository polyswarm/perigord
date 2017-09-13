// Invokes the perigord driver application

package main

import (
	_ "{{.project}}/migrations"
	"github.com/swarmdotmarket/perigord/stub"
)

func main() {
	stub.StubMain()
}
