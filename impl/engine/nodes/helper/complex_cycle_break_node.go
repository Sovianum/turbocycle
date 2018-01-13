package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type ComplexCycleBreakNode interface {
	graph.Node
	nodes.ComplexGasChannel
}

func NewComplexCycleBreakNode(gas gases.Gas, tStag, pStag, massRate float64) ComplexCycleBreakNode {
	var result = &complexCycleBreakNode{}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
			&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
		},
		[]string{
			nodes.TemperatureInputTag, nodes.PressureInputTag, nodes.GasInputTag, nodes.MassRateInputTag,
			nodes.TemperatureOutputTag, nodes.PressureOutputTag, nodes.GasOutputTag, nodes.MassRateOutputTag,
		},
	)

	graph.SetAll(
		[]graph.PortState{
			states.NewGasPortState(gas), states.NewTemperaturePortState(tStag),
			states.NewPressurePortState(pStag), states.NewMassRatePortState(massRate),
		},
		[]graph.Port{
			result.gasInput, result.temperatureInput, result.pressureInput, result.massRateInput,
		},
	)

	return result
}

type complexCycleBreakNode struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port
}

func (node *complexCycleBreakNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *complexCycleBreakNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *complexCycleBreakNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *complexCycleBreakNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *complexCycleBreakNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *complexCycleBreakNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *complexCycleBreakNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *complexCycleBreakNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *complexCycleBreakNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ComplexCycleBreak")
}

func (node *complexCycleBreakNode) GetRequirePorts() []graph.Port {
	return make([]graph.Port, 0)
}

func (node *complexCycleBreakNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.temperatureOutput,
		node.pressureOutput,
		node.gasOutput,
		node.massRateOutput,
	}
}

func (node *complexCycleBreakNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput,
		node.pressureInput,
		node.gasInput,
		node.massRateInput,

		node.temperatureOutput,
		node.pressureOutput,
		node.gasOutput,
		node.massRateOutput,
	}
}

func (node *complexCycleBreakNode) Process() error {
	graph.SetAll(
		[]graph.PortState{
			node.temperatureInput.GetState(),
			node.pressureInput.GetState(),
			node.gasInput.GetState(),
			node.massRateInput.GetState(),
		},
		[]graph.Port{
			node.temperatureOutput,
			node.pressureOutput,
			node.gasOutput,
			node.massRateOutput,
		},
	)
	return nil
}
