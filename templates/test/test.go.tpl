package tests

import (
	"gopkg.in/check.v1"

	"github.com/polyswarm/perigord/contract"
	"github.com/polyswarm/perigord/network"
	"github.com/polyswarm/perigord/testing"

	"{{.project}}/bindings"
)

type {{.test}}Suite struct {
    network     *network.Network
}

var _ = check.Suite(&{{.test}}Suite{})

func (s *{{.test}}Suite) SetUpTest(c *check.C) {
	nw, err := testing.SetUpTest()
	if err != nil {
		c.Fatal(err)
	}

	s.network = nw
}

func (s *{{.test}}Suite) TearDownTest(c *check.C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE
