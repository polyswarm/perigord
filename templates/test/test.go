package tests

import (
	. "gopkg.in/check.v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/swarmdotmarket/perigord/contract"
	"github.com/swarmdotmarket/perigord/testing"

	"github.com/swarmdotmarket/token/bindings"
)

type {{.test}}Suite struct {
	auth    *bind.TransactOpts
	backend bind.ContractBackend
}

var _ = Suite(&{{.test}}Suite{})

func (s *{{.test}}Suite) SetUpTest(c *C) {
	auth, backend := testing.SetUpTest()

	s.auth = auth
	s.backend = backend
}

func (s *{{.test}}Suite) TearDownTest(c *C) {
	testing.TearDownTest()
}

// USER TESTS GO HERE
