package spy

import (
	"github.com/super-l/machine-code/machine"
	"runtime"
)

type registerPack struct {
	App     string `json:"app,omitempty"`
	Version string `json:"version,omitempty"`
	CPUID   string `json:"CPUID,omitempty"`
	UUID    string `json:"UUID,omitempty"`
	SN      string `json:"SN,omitempty"`
	GO      string `json:"GO,omitempty"`
}

var regPack registerPack

func init() {
	regPack.CPUID, _ = machine.GetCpuId()
	regPack.UUID, _ = machine.GetPlatformUUID()
	regPack.SN, _ = machine.GetSerialNumber()
	regPack.GO = runtime.Version()
}

func Register(app, version string) {
	regPack.App = app
	regPack.Version = version
}
