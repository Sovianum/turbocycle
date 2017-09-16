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
	portA         = "portA"
	portB         = "portB"
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

type GasSource interface {
	GasOutput() core.Port
	TStagOut() float64
	PStagOut() float64
}

type GasSink interface {
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
