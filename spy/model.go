package spy

import (
	"github.com/super-l/machine-code/machine"
	"runtime"
)

type Register struct {
	App     string `json:"app,omitempty"`
	Version string `json:"version,omitempty"`
	CPUID   string `json:"CPUID,omitempty"`
	UUID    string `json:"UUID,omitempty"`
	SN      string `json:"SN,omitempty"`
	GO      string `json:"GO,omitempty"`
}

var RegisterPackage *Register

func init() {
	RegisterPackage.CPUID, _ = machine.GetCpuId()
	RegisterPackage.UUID, _ = machine.GetPlatformUUID()
	RegisterPackage.SN, _ = machine.GetSerialNumber()
	RegisterPackage.GO = runtime.Version()
}
