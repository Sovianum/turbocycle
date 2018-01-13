package constructive

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewBaseBlockedTurbine(node graph.Node, precision float64) *baseBlockedTurbine {
	var result = &baseBlockedTurbine{
		precision: precision,
	}
	graph.AttachAllWithTags(
		node,
		[]*graph.Port{
			&result.powerOutput,
			&result.temperatureInput, &result.pressureInput, &result.gasInput,
			&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
		},
		[]string{
			nodes.PowerOutputTag,
			nodes.TemperatureInputTag, nodes.PressureInputTag, nodes.GasInputTag,
			nodes.TemperatureOutputTag, nodes.PressureOutputTag, nodes.GasOutputTag, nodes.MassRateOutputTag,
		},
	)
	return result
}

type baseBlockedTurbine struct {
	graph.BaseNode

	powerOutput graph.Port

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	precision float64
}

func (node *baseBlockedTurbine) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureInput, node.pressureInput, node.gasInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBlockedTurbine) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput,
	}
}

func (node *baseBlockedTurbine) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *baseBlockedTurbine) MassRateRel() float64 {
	return node.massRateOutput.GetState().(states.MassRatePortState).MassRate
}

func (node *baseBlockedTurbine) InputGas() gases.Gas {
	return node.inputGas()
}

func (node *baseBlockedTurbine) TStagIn() float64 {
	return node.tStagIn()
}

func (node *baseBlockedTurbine) PStagIn() float64 {
	return node.pStagIn()
}

func (node *baseBlockedTurbine) TStagOut() float64 {
	return node.tStagOut()
}

func (node *baseBlockedTurbine) PStagOut() float64 {
	return node.pStagOut()
}

func (node *baseBlockedTurbine) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *baseBlockedTurbine) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *baseBlockedTurbine) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *baseBlockedTurbine) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *baseBlockedTurbine) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *baseBlockedTurbine) GasInput() graph.Port {
	return node.gasInput
}

func (node *baseBlockedTurbine) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *baseBlockedTurbine) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *baseBlockedTurbine) inputGas() gases.Gas {
	return node.gasInput.GetState().(states.GasPortState).Gas
}

func (node *baseBlockedTurbine) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBlockedTurbine) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *baseBlockedTurbine) tStagOut() float64 {
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseBlockedTurbine) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}
