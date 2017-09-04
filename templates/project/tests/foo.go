package tests

import (
	. "gopkg.in/check.v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/swarmdotmarket/perigord/testing"

	"{{.project}}/bindings"
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
	_, _, interactor, _ := bindings.DeployFoo(testing.Auth(), testing.Backend())
	testing.Backend().Commit()

	ret, _ := interactor.Bar(&bind.CallOpts{Pending: true})
	c.Assert(int64(1337), Equals, ret.Int64())
}
