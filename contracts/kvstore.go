package contracts

import "vastchain.ltd/vastchain/chain_structure"

type kvStoreContract struct{}

var KVStoreContract = NewBuiltinContract("KVStore", kvStoreContract{})

func (token *kvStoreContract) C(args []chain_structure.VcContractTypedValue) ([]chain_structure.VcContractTypedValue, error) {
	return nil, nil // chain_structure.NewVcErrorNotImplemented()
}
