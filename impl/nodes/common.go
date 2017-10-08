package nodes

import "github.com/Sovianum/turbocycle/core"

const (
	PressureInput              = "pressureInput"
	PressureOutput             = "pressureOutput"
	TemperatureInput           = "temperatureInput"
	TemperatureOutput          = "temperatureOutput"
	GasInput                   = "gasInput"
	GasOutput                  = "gasOutput"
	ComplexGasPort             = "complexGasPort"
	PressurePort               = "pressurePort"
	TemperaturePort            = "temperaturePort"
	MassRateRelPort            = "massRateRelPort"
	MassRateRelInput           = "massRateRelInput"
	MassRateRelOutput          = "massRateRelOutput"
	GasPort                    = "gasPort"
	PowerInput                 = "powerInput"
	PowerOutput                = "powerOutput"
	ComplexGasInput            = "complexGasInput"
	ComplexGasOutput           = "complexGasOutput"
	ColdGasInput               = "coldGasInput"
	ColdGasOutput              = "coldGasOutput"
	HotGasInput                = "hotGasInput"
	HotGasOutput               = "hotGasOutput"
	UpdatePort                 = "updatePort"
	DataSourcePort             = "dataSourcePort"
	CompressorComplexGasInput  = "compressorComplexGasInput"
	CompressorComplexGasOutput = "compressorComplexGasOutput"
	TurbineComplexGasInput     = "turbineComplexGasInput"
	TurbineComplexGasOutput    = "turbineComplexGasOutput"
	DefaultN                   = 50
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
