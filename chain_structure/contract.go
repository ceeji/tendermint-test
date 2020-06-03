package chain_structure

type Contract interface {
	Call(name string, args []VcContractTypedValue) ([]VcContractTypedValue, error)
}

type ContractChecker interface {
	Check(name string, args []VcContractTypedValue) error
}
