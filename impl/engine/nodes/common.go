package nodes

import "github.com/Sovianum/turbocycle/core/graph"

const (
	DefaultN = 50
)

type ComplexGasChannel interface {
	ComplexGasSink
	ComplexGasSource
}

type ComplexGasSource interface {
	ComplexGasOutput() graph.Port
}

type ComplexGasSink interface {
	ComplexGasInput() graph.Port
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

type MassRateRelChannel interface {
	MassRateRelSource
	MassRateRelSink
}

type MassRateRelSource interface {
	MassRateRelOutput() graph.Port
	MassRateRelOut() float64
}

type MassRateRelSink interface {
	MassRateRelInput() graph.Port
	MassRateRelIn() float64
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

type TemperatureOut interface {
	TStagOut() float64
}

type TemperatureSink interface {
	TemperatureInput() graph.Port
}

type TemperatureIn interface {
	TStagIn() float64
}

type PressureChannel interface {
	PressureSource
	PressureSink
}

type PressureSource interface {
	PressureOutput() graph.Port
}

type PressureOut interface {
	PStagOut() float64
}

type PressureSink interface {
	PressureInput() graph.Port
}

type PressureIn interface {
	PStagIn() float64
}

func IsDataSource(port graph.Port) (bool, error) {
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
