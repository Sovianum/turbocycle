package nodes

import "github.com/Sovianum/turbocycle/core"

const (
	DefaultN = 50
)

type ComplexGasChannel interface {
	ComplexGasSink
	ComplexGasSource
}

type ComplexGasSource interface {
	ComplexGasOutput() core.Port
}

type ComplexGasSink interface {
	ComplexGasInput() core.Port
}

type PowerChannel interface {
	PowerSink
	PowerSource
}

type PowerSource interface {
	PowerOutput() core.Port
}

type PowerSink interface {
	PowerInput() core.Port
}

type MassRateRelChannel interface {
	MassRateRelSource
	MassRateRelSink
}

type MassRateRelSource interface {
	MassRateRelOutput() core.Port
	MassRateRelOut() float64
}

type MassRateRelSink interface {
	MassRateRelInput() core.Port
	MassRateRelIn() float64
}

type GasChannel interface {
	GasSource
	GasSink
}

type GasSource interface {
	GasOutput() core.Port
}

type GasSink interface {
	GasInput() core.Port
}

type TemperatureChannel interface {
	TemperatureSource
	TemperatureSink
}

type TemperatureSource interface {
	TemperatureOutput() core.Port
}

type TemperatureOut interface {
	TStagOut() float64
}

type TemperatureSink interface {
	TemperatureInput() core.Port
}

type TemperatureIn interface {
	TStagIn() float64
}

type PressureChannel interface {
	PressureSource
	PressureSink
}

type PressureSource interface {
	PressureOutput() core.Port
}

type PressureOut interface {
	PStagOut() float64
}

type PressureSink interface {
	PressureInput() core.Port
}

type PressureIn interface {
	PStagIn() float64
}

func IsDataSource(port core.Port) (bool, error) {
	var linkPort = port.GetLinkPort()
	if linkPort == nil {
		return false, nil
	}

	var outerNode = port.GetOuterNode()
	if outerNode == nil {
		return false, nil
	}

	if !outerNode.ContextDefined() {
		return false, nil
	}

	var updatePorts = outerNode.GetUpdatePorts()

	for _, port := range updatePorts {
		if port == linkPort {
			return true, nil
		}
	}

	return false, nil
}
