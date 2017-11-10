package tests

import (
	. "gopkg.in/check.v1"

	"github.com/polyswarm/perigord/contract"
	"github.com/polyswarm/perigord/testing"

	"{{.project}}/bindings"
)

type {{.test}}Suite struct {}

var _ = Suite(&{{.test}}Suite{})

func (s *{{.test}}Suite) SetUpTest(c *C) {
	testing.SetUpTest()
}

func (s *{{.test}}Suite) TearDownTest(c *C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE
