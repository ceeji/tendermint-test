package chain_structure

import (
	"bytes"
	"text/template"
	"vastchain.ltd/vastchain/utils"
)

type VcError struct {
	Code                string // format: VC12503
	Name                string
	descriptionTemplate *template.Template // support go template to show instance-specific message
	Detail              map[string]string  // more information

}

var errorDescTemplateCache = make(map[string]*template.Template)
var emptyDetail = make(map[string]string)

// NewVcError creates a VcError instance
func NewVcError(code, name, descTemplate string, detail map[string]string) *VcError {
	cachedTemplate, ok := errorDescTemplateCache[code]
	if !ok {
		var err error
		cachedTemplate, err = template.New("error").Parse(descTemplate)
		if err != nil {
			cachedTemplate, _ = template.New("error").Parse("!error compiling error description template")
		}

		errorDescTemplateCache[code] = cachedTemplate
	}

	if detail == nil {
		detail = emptyDetail
	}

	return &VcError{
		Code:                code,
		Name:                name,
		descriptionTemplate: cachedTemplate,
		Detail:              detail,
	}
}

func (error *VcError) Error() string {
	return "[VC" + error.Code + "] " + error.Name + ": " + error.Description()
}

func (error *VcError) Description() string {
	writer := new(bytes.Buffer)
	_ = error.descriptionTemplate.Execute(writer, error)
	desc := utils.BytesToString(writer.Bytes())

	return desc
}
