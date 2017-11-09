// Invokes the perigord driver application

package main

import (
	_ "{{.project}}/migrations"
	"github.com/polyswarm/perigord/stub"
)

func main() {
	stub.StubMain()
}
