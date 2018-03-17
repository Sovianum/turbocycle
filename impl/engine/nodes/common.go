package nodes

import (
	"math"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/material/gases"
)

const (
	DefaultN             = 50
	RPMInputTag          = "rpmInput"
	RPMOutputTag         = "rpmOutput"
	PowerInputTag        = "powerInput"
	PowerOutputTag       = "powerOutput"
	GasInputTag          = "gasInput"
	GasOutputTag         = "gasOutput"
	PressureInputTag     = "pressureInput"
	PressureOutputTag    = "pressureOutput"
	TemperatureInputTag  = "temperatureInput"
	TemperatureOutputTag = "temperatureOutput"
	MassRateInputTag     = "massRateInput"
	MassRateOutputTag    = "massRateOutput"
)

type EnthalpyChannel interface {
	TemperatureChannel
	MassRateChannel
	GasChannel
}

func EnthalpyDiff(ch EnthalpyChannel) float64 {
	mri := ch.MassRateInput().GetState().Value().(float64)
	mro := ch.MassRateOutput().GetState().Value().(float64)

	iGas := ch.GasInput().GetState().Value().(gases.Gas)
	oGas := ch.GasOutput().GetState().Value().(gases.Gas)

	it := ch.TemperatureInput().GetState().Value().(float64)
	ot := ch.TemperatureOutput().GetState().Value().(float64)

	tMin := math.Min(it, ot)

	return mro*(ot-tMin)*gases.CpMean(oGas, tMin, ot, DefaultN) - mri*(it-tMin)*gases.CpMean(iGas, tMin, it, DefaultN)
}

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

	if !outerNode.ContextDefined(0) {
		return false, nil
	}

	var updatePorts, _ = outerNode.GetUpdatePorts()

	for _, port := range updatePorts {
		if port == linkPort {
			return true, nil
		}
	}

	return false, nil
}

func LinkComplexInToOut(node1 ComplexGasSink, node2 ComplexGasSource) {
	LinkComplexOutToIn(node2, node1)
}

func LinkComplexOutToIn(node1 ComplexGasSource, node2 ComplexGasSink) {
	graph.LinkAll(
		[]graph.Port{node1.GasOutput(), node1.TemperatureOutput(), node1.PressureOutput(), node1.MassRateOutput()},
		[]graph.Port{node2.GasInput(), node2.TemperatureInput(), node2.PressureInput(), node2.MassRateInput()},
	)
}

func LinkComplexOutToOut(node1 ComplexGasSource, node2 ComplexGasSource) {
	graph.LinkAll(
		[]graph.Port{node1.GasOutput(), node1.TemperatureOutput(), node1.PressureOutput(), node1.MassRateOutput()},
		[]graph.Port{node2.GasOutput(), node2.TemperatureOutput(), node2.PressureOutput(), node2.MassRateOutput()},
	)
}

func LinkComplexInToIn(node1 ComplexGasSink, node2 ComplexGasSink) {
	graph.LinkAll(
		[]graph.Port{node1.GasInput(), node1.TemperatureInput(), node1.PressureInput(), node1.MassRateInput()},
		[]graph.Port{node2.GasInput(), node2.TemperatureInput(), node2.PressureInput(), node2.MassRateInput()},
	)
}
