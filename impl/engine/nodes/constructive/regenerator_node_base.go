package constructive

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func newBaseRegenerator(node graph.Node) *baseRegenerator {
	var result = &baseRegenerator{}

	graph.AttachAllWithTags(
		node,
		[]*graph.Port{
			&result.hotTemperatureInput, &result.hotPressureInput, &result.hotGasInput, &result.hotMassRateInput,
			&result.hotTemperatureOutput, &result.hotPressureOutput, &result.hotGasOutput, &result.hotMassRateOutput,
			&result.coldTemperatureInput, &result.coldPressureInput, &result.coldGasInput, &result.coldMassRateInput,
			&result.coldTemperatureOutput, &result.coldPressureOutput, &result.coldGasOutput, &result.coldMassRateOutput,
		},
		[]string{
			"hotTemperatureInput", "hotPressureInput", "hotGasInput", "hotMassRateInput",
			"hotTemperatureOutput", "hotPressureOutput", "hotGasOutput", "hotMassRateOutput",
			"coldTemperatureInput", "coldPressureInput", "coldGasInput", "coldMassRateInput",
			"coldTemperatureOutput", "coldPressureOutput", "coldGasOutput", "coldMassRateOutput",
		},
	)
	return result
}

type baseRegenerator struct {
	graph.BaseNode

	hotTemperatureInput graph.Port
	hotPressureInput    graph.Port
	hotGasInput         graph.Port
	hotMassRateInput    graph.Port

	hotTemperatureOutput graph.Port
	hotPressureOutput    graph.Port
	hotGasOutput         graph.Port
	hotMassRateOutput    graph.Port

	coldTemperatureInput graph.Port
	coldPressureInput    graph.Port
	coldGasInput         graph.Port
	coldMassRateInput    graph.Port

	coldTemperatureOutput graph.Port
	coldPressureOutput    graph.Port
	coldGasOutput         graph.Port
	coldMassRateOutput    graph.Port
}

func (node *baseRegenerator) HotInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.hotGasInput, node.hotTemperatureInput, node.hotPressureInput, node.hotMassRateInput,
	)
}

func (node *baseRegenerator) HotOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.hotGasOutput, node.hotTemperatureOutput, node.hotPressureOutput, node.hotMassRateOutput,
	)
}

func (node *baseRegenerator) ColdInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.coldGasInput, node.coldTemperatureInput, node.coldPressureInput, node.coldMassRateInput,
	)
}

func (node *baseRegenerator) ColdOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.coldGasOutput, node.coldTemperatureOutput, node.coldPressureOutput, node.coldMassRateOutput,
	)
}

func (node *baseRegenerator) GetPorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureInput, node.hotPressureInput, node.hotGasInput, node.hotMassRateInput,
		node.hotTemperatureOutput, node.hotPressureOutput, node.hotGasOutput, node.hotMassRateOutput,
		node.coldTemperatureInput, node.coldPressureInput, node.coldGasInput, node.coldMassRateInput,
		node.coldTemperatureOutput, node.coldPressureOutput, node.coldGasOutput, node.coldMassRateOutput,
	}
}

func (node *baseRegenerator) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureInput, node.hotPressureInput, node.hotGasInput, node.hotMassRateInput,
		node.coldTemperatureInput, node.coldPressureInput, node.coldGasInput, node.coldMassRateInput,
	}
}

func (node *baseRegenerator) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.hotTemperatureOutput, node.hotPressureOutput, node.hotGasOutput, node.hotMassRateOutput,
		node.coldTemperatureOutput, node.coldPressureOutput, node.coldGasOutput, node.coldMassRateOutput,
	}
}

func (node *baseRegenerator) tStagHotIn() float64 {
	return node.hotTemperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseRegenerator) tStagHotOut() float64 {
	return node.hotTemperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseRegenerator) tStagColdIn() float64 {
	return node.coldTemperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *baseRegenerator) tStagColdOut() float64 {
	return node.coldTemperatureOutput.GetState().(states.TemperaturePortState).TStag
}
