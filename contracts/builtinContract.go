package contracts

import (
	"fmt"
	"reflect"
	"vastchain.ltd/vastchain/chain_structure"
)

type BuiltinContract struct {
	Name             string
	impl             interface{}
	implPointerValue reflect.Value
	implType         reflect.Type
	methodTable      map[string]*reflect.Value
}

func NewBuiltinContract(name string, impl interface{}) *BuiltinContract {
	ret := &BuiltinContract{
		Name:             name,
		impl:             impl,
		implType:         reflect.PtrTo(reflect.ValueOf(impl).Type()),
		methodTable:      make(map[string]*reflect.Value),
		implPointerValue: reflect.New(reflect.ValueOf(impl).Type()),
	}

	ret.implPointerValue.Elem().Set(reflect.ValueOf(impl))
	typeName := ret.implType.Kind()
	fmt.Print(typeName)
	methodCount := ret.implType.NumMethod()
	for i := 0; i < methodCount; i++ {
		method := ret.implType.Method(i)
		ret.methodTable[method.Name] = &method.Func
	}

	return ret
}

func (contract *BuiltinContract) Call(name string, args []chain_structure.VcContractTypedValue) ([]chain_structure.VcContractTypedValue, error) {
	method := contract.methodTable[name]
	if method == nil {
		return nil, chain_structure.NewVcErrorContractNoSuchMethodError(name, contract.Name)
	}

	ret := method.Call([]reflect.Value{contract.implPointerValue, reflect.ValueOf(args)})
	result := ret[0].Interface().([]chain_structure.VcContractTypedValue)
	errInterface := ret[1]

	if errInterface.IsNil() {
		return result, nil
	}
	return result, errInterface.Interface().(error)
}
