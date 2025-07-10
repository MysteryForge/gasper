package integrity

import (
	"go.k6.io/k6/js/modules"
)

func init() {}

type RootModule struct{}

type ModuleInstance struct {
	vu modules.VU
}

var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

func New() *RootModule {
	return &RootModule{}
}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
	}
}

func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"sayHello": func() interface{} {
				return "Hello Integrity"
			},
		},
	}
}
