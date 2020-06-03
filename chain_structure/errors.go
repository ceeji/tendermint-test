package chain_structure

// NewVcErrorInvalidArgument creates InvalidArgument error
func NewVcErrorInvalidArgument(argumentName string, reason string) *VcError {
	return NewVcError("10001", "invalid_argument",
		"the provided argument {{if .Detail.name}}\"{{.Detail.name}}\" {{end}}is invalid{{if .Detail.reason}}: {{.Detail.reason}}{{end}}",
		map[string]string{
			"name":   argumentName,
			"reason": reason,
		})
}

// NewVcErrorInvalidArgument creates InvalidArgument error
func NewVcErrorPreconditionNotSatisfied(precondition string) *VcError {
	return NewVcError("10002", "precondition_not_satisfied",
		"the precondition is not satisfied{{if .Detail.precondition}}: {{.Detail.precondition}}{{end}}",
		map[string]string{
			"precondition": precondition,
		})
}

// NewVcErrorContractNoSuchMethodError creates VcErrorContractNoSuchMethod error
func NewVcErrorContractNoSuchMethodError(method, contract string) *VcError {
	return NewVcError("20001", "contract_no_such_method",
		"the provided method {{.Detail.name}} not found in contract {{.Detail.contract}}",
		map[string]string{
			"name":     method,
			"contract": contract,
		})
}

// NewNotImplementedError creates NotImplemented error
func NewVcErrorNotImplemented() *VcError {
	return NewVcError("90001", "not_implemented",
		"the function is not implemented", nil)
}

// NewVcErrorNotInitialized creates NotInitialized error
func NewVcErrorNotInitialized() *VcError {
	return NewVcError("90002", "not_initialized",
		"the function is not initialized", nil)
}
