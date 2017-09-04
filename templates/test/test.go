package tests

import (
	. "gopkg.in/check.v1"

	"github.com/swarmdotmarket/perigord/testing"

	"{{.project}}/bindings"
)

type {{.test}} struct{}

var _ = Suite(&{{.test}}{})

func (s *{{.test}}) SetUpTest(c *C) {
	testing.SetUpTest()
}

func (s *{{.test}}) TearDownTest(c *C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE

func (s *{{.test}}) TestDummy(c *C) {
	c.Succeed()
}
