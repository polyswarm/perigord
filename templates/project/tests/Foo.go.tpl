package tests

import (
	"gopkg.in/check.v1"

	"github.com/polyswarm/perigord/contract"
	"github.com/polyswarm/perigord/testing"

	"{{.project}}/bindings"
)

type foo_test struct{}

var _ = Suite(&foo_test{})

func (s *foo_test) SetUpTest(c *check.C) {
	testing.SetUpTest()
}

func (s *foo_test) TearDownTest(c *check.C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE

func (s *foo_test) TestFoo(c *check.C) {
	session := contract.Session("Foo")
	c.Assert(session, NotNil)

	foo_session, ok := session.(*bindings.FooSession)
	c.Assert(ok, Equals, true)
	c.Assert(foo_session, NotNil)

	ret, _ := foo_session.Bar()
	c.Assert(int64(1337), Equals, ret.Int64())
}
