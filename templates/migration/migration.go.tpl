package migrations

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/polyswarm/perigord/contract"
	"github.com/polyswarm/perigord/migration"

	"{{.project}}/bindings"
)

type {{.contract}}Deployer struct{}

func (d *{{.contract}}Deployer) Deploy(ctx context.Context, network *migration.Network) (common.Address, *types.Transaction, interface{}, error) {
    auth := network.NewTransactor(0)
	address, transaction, contract, err := bindings.Deploy{{.contract}}(auth, network.Backend())
	if err != nil {
		return common.Address{}, nil, nil, err
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

func (d *{{.contract}}Deployer) Bind(ctx context.Context, network *migration.Network, address common.Address) (interface{}, error) {
    auth := network.NewTransactor(0)
	contract, err := bindings.New{{.contract}}(address, network.Backend())
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
		F: func(ctx context.Context, network *migration.Network) error {
			if err := contract.Deploy(ctx, "{{.contract}}", network); err != nil {
				return err
			}

			return nil
		},
	})
}
