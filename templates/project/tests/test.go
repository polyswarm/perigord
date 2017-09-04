package tests

import (
	. "gopkg.in/check.v1"

	"github.com/swarmdotmarket/perigord/testing"
)

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
	testing.SetUpTest()
}

func (s *suite) TearDownTest(c *C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE

func (s *suite) TestFoo(c *C) {
	c.Succeed()
}
