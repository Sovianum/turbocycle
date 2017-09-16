package nodes

import "github.com/Sovianum/turbocycle/core"

const (
	pressureInput     = "pressureInput"
	pressureOutput    = "pressureOutput"
	temperatureInput  = "temperatureInput"
	temperatureOutput = "temperatureOutput"
	gasInput          = "gasInput"
	gasOutput         = "gasOutput"
	complexGasPort    = "complexGasPort"
	pressurePort      = "pressurePort"
	temperaturePort   = "temperaturePort"
	massRateRelPort   = "massRateRelPort"
	gasPort           = "gasPort"
	powerInput        = "powerInput"
	powerOutput       = "powerOutput"
	complexGasInput   = "complexGasInput"
	complexGasOutput  = "complexGasOutput"
	coldGasInput      = "coldGasInput"
	coldGasOutput     = "coldGasOutput"
	hotGasInput       = "hotGasInput"
	hotGasOutput      = "hotGasOutput"
	portA             = "portA"
	portB             = "portB"
	defaultN          = 50
)

type ComplexGasChannel interface {
	ComplexGasSink
	ComplexGasSource
}

type ComplexGasSource interface {
	ComplexGasOutput() core.Port
	TStagOut() float64
	PStagOut() float64
}

type ComplexGasSink interface {
	ComplexGasInput() core.Port
	TStagIn() float64
	PStagIn() float64
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
	TStag() float64
}

type TemperatureSink interface {
	TemperatureInput() core.Port
	TStagIn() float64
}

type PressureChannel interface {
	PressureSource
	PressureSink
}

type PressureSource interface {
	PressureOutput() core.Port
	PStagOut() float64
}

type PressureSink interface {
	PressureInput() core.Port
	PStagIn() float64
}

func isDataSource(port core.Port) (bool, error) {
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

	var updatePortTags, err = outerNode.GetUpdatePortTags()
	if err != nil {
		return false, err
	}

	for _, updatePortTag := range updatePortTags {
		var tagPort, tagErr = outerNode.GetPortByTag(updatePortTag)
		if tagErr != nil {
			return false, tagErr
		}

		if tagPort == linkPort {
			return true, nil
		}
	}

	return false, nil
}
