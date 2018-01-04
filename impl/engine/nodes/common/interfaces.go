package common

import "github.com/Sovianum/turbocycle/core/graph"

type RPMChannel interface {
	RPMSink
	RPMSource
}

type RPMSource interface {
	RPMOutput() graph.Port
}

type RPMSink interface {
	RPMInput() graph.Port
}

type PowerChannel interface {
	PowerSink
	PowerSource
}

type PowerSource interface {
	PowerOutput() graph.Port
}

type PowerSink interface {
	PowerInput() graph.Port
}

type ComplexGasChannel interface {
	ComplexGasSink
	ComplexGasSource
}

type ComplexGasSource interface {
	GasSource
	TemperatureSource
	PressureSource
	MassRateSource
}

type ComplexGasSink interface {
	GasSink
	TemperatureSink
	PressureSink
	MassRateSink
}

type MassRateChannel interface {
	MassRateInput() graph.Port
	MassRateOutput() graph.Port
}

type MassRateSource interface {
	MassRateOutput() graph.Port
}

type MassRateSink interface {
	MassRateInput() graph.Port
}

type GasChannel interface {
	GasSource
	GasSink
}

type GasSource interface {
	GasOutput() graph.Port
}

type GasSink interface {
	GasInput() graph.Port
}

type TemperatureChannel interface {
	TemperatureSource
	TemperatureSink
}

type TemperatureSource interface {
	TemperatureOutput() graph.Port
}

type TemperatureSink interface {
	TemperatureInput() graph.Port
}

type PressureChannel interface {
	PressureSource
	PressureSink
}

type PressureSource interface {
	PressureOutput() graph.Port
}

type PressureSink interface {
	PressureInput() graph.Port
}

type TemperatureOut interface {
	TStagOut() float64
}

type TemperatureIn interface {
	TStagIn() float64
}

type PressureOut interface {
	PStagOut() float64
}

type PressureIn interface {
	PStagIn() float64
}
