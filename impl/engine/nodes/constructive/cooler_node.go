package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type CoolerNode interface {
	graph.Node

	nodes.GasChannel
	nodes.PressureChannel
	nodes.TemperatureChannel
	nodes.MassRateChannel
}

func NewCoolerNode(tOut float64, sigma float64) CoolerNode {
	var result = &coolerNode{
		tOut:  tOut,
		sigma: sigma,
	}

	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.gasInput, &result.temperatureInput, &result.pressureInput, &result.massRateInput,
			&result.gasOutput, &result.temperatureOutput, &result.pressureOutput, &result.massRateOutput,
		},
		[]string{
			nodes.GasInputTag, nodes.TemperatureInputTag, nodes.PressureInputTag, nodes.MassRateInputTag,
			nodes.GasOutputTag, nodes.TemperatureOutputTag, nodes.PowerOutputTag, nodes.MassRateOutputTag,
		},
	)

	return result
}

type coolerNode struct {
	graph.BaseNode

	gasInput         graph.Port
	temperatureInput graph.Port
	pressureInput    graph.Port
	massRateInput    graph.Port

	gasOutput         graph.Port
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	massRateOutput    graph.Port

	tOut  float64
	sigma float64
}

func (node *coolerNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *coolerNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *coolerNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *coolerNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *coolerNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *coolerNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *coolerNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *coolerNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *coolerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Cooler")
}

func (node *coolerNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput,
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *coolerNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput,
	}
}

func (node *coolerNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	}
}

func (node *coolerNode) Process() error {
	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(),
			states.NewTemperaturePortState(node.tOut),
			states.NewPressurePortState(node.pressureInput.GetState().(states.PressurePortState).PStag * node.sigma),
			node.massRateInput.GetState(),
		},
		[]graph.Port{node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput},
	)
	return nil
}
