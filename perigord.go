// Copyright © 2017 PolySwarm <info@polyswarm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Misc utility functions for talking with go-ethereum

package perigord

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func CheckOutOfGas(tx *types.Transaction, rcpt *types.Receipt) bool {
	return tx.Gas().Cmp(rcpt.GasUsed) == 0
}

func WaitMined(ctx context.Context, backend bind.DeployBackend, tx *types.Transaction) (*types.Receipt, error) {
	rcpt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		return nil, err
	}
	if CheckOutOfGas(tx, rcpt) {
		return nil, errors.New("out of gas")
	}

	return rcpt, nil
}

func EventSignatureToTopicHash(signature string) common.Hash {
	return crypto.Keccak256Hash([]byte(signature))
}
