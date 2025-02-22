package js

import "github.com/dop251/goja"

func NewRuntime() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	return vm
}
