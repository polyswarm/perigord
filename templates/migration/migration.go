package migrations

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/swarmdotmarket/perigord/contract"
	"github.com/swarmdotmarket/perigord/migration"

	"{{.project}}/bindings"
)

type {{.contract}}Deployer struct{}

func (d *{{.contract}}Deployer) Deploy(ctx context.Context, auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, interface{}, error) {
	address, transaction, contract, err := bindings.Deploy{{.contract}}(auth, backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	val, ok := backend.(*backends.SimulatedBackend)
	if ok {
		val.Commit()
	}

	session := &bindings.{{.contract}}Session{
		Contract: contract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: *auth,
	}

	return address, transaction, session, nil
}

func (d *{{.contract}}Deployer) Bind(ctx context.Context, auth *bind.TransactOpts, backend bind.ContractBackend, address common.Address) (interface{}, error) {
	contract, err := bindings.New{{.contract}}(address, backend)
	if err != nil {
		return nil, err
	}

	session := &bindings.{{.contract}}Session{
		Contract: contract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: *auth,
	}

	return session, nil
}

func init() {
	contract.AddContract("{{.contract}}", &{{.contract}}Deployer{})

	migration.AddMigration(&migration.Migration{
		Number: {{.number}},
		F: func(ctx context.Context, auth *bind.TransactOpts, backend bind.ContractBackend) error {
			if err := contract.Deploy(ctx, "{{.contract}}", auth, backend); err != nil {
				return err
			}

			return nil
		},
	})
}
