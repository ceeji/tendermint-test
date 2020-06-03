package chain_structure

// VcAction represents action in VastChain
type VcAction struct {
	Contract string                 // name of the contract to be called
	Function string                 // name of the function
	Args     []VcContractTypedValue // args
}

// VcContractTypedValue represents a value with a type prefix
type VcContractTypedValue struct {
	Type  byte
	Value interface{}
}

// TODO: add function to check the value of VcContractTypedValue

type VcContractDataType byte

const (
	VcContractDataTypeInt64 VcContractDataType = iota
	VcContractDataTypeInt64Array
	VcContractDataTypeString
	VcContractDataTypeStringArray
	VcContractDataTypeByteArray
	VcContractDataTypeAddress
	VcContractDataTypeAddressArray
	VcContractDataTypeTypedValueArray
)
