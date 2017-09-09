package nodes

import "github.com/Sovianum/turbocycle/core"

const (
	powerInput    = "powerInput"
	powerOutput   = "powerOutput"
	gasInput      = "gasInput"
	gasOutput     = "gasOutput"
	coldGasInput  = "coldGasInput"
	coldGasOutput = "coldGasOutput"
	hotGasInput   = "hotGasInput"
	hotGasOutput  = "hotGasOutput"
	defaultN      = 50
)

type GasChannel interface {
	GasSink
	GasSource
}

type PowerChannel interface {
	PowerSink
	PowerSource
}

type GasSink interface {
	GasOutput() core.Port
	TStagOut() float64
	PStagOut() float64
}

type GasSource interface {
	GasInput() core.Port
	TStagIn() float64
	PStagIn() float64
}

type PowerSource interface {
	PowerOutput() core.Port
}

type PowerSink interface {
	PowerInput() core.Port
}
