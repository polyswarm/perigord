package tests

import (
	. "gopkg.in/check.v1"

	"github.com/swarmdotmarket/perigord/contract"
	"github.com/swarmdotmarket/perigord/testing"

	"github.com/swarmdotmarket/token/bindings"
)

type foo_test struct{}

var _ = Suite(&foo_test{})

func (s *foo_test) SetUpTest(c *C) {
	testing.SetUpTest()
}

func (s *foo_test) TearDownTest(c *C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE

func (s *foo_test) TestFoo(c *C) {
	session := contract.Session("Foo")
	c.Assert(session, NotNil)

	foo_session, ok := session.(*bindings.FooSession)
	c.Assert(ok, Equals, true)
	c.Assert(foo_session, NotNil)

	ret, _ := foo_session.Bar()
	c.Assert(int64(1337), Equals, ret.Int64())
}
