package contracts

import (
	"vastchain.ltd/vastchain/chain_structure"
)

type tokenContract struct{}

var TokenContract = NewBuiltinContract("Token", tokenContract{})

func (token *tokenContract) Create(args []chain_structure.VcContractTypedValue) ([]chain_structure.VcContractTypedValue, error) {
	return nil, nil // chain_structure.NewVcErrorNotImplemented()
}
