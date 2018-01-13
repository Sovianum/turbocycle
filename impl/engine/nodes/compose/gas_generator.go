package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
)

type GasGeneratorNode interface {
	graph.Node
	nodes.ComplexGasChannel
	Burner() constructive.BurnerNode
	TurboCascade() TurboCascadeNode
}

func NewGasGeneratorNode(
	compressorEtaAd, piStag float64,
	fuel fuel.GasFuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0 float64,
	etaT, lambdaOut float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) GasGeneratorNode {
	var result = &gasGeneratorNode{
		turboCascade: NewTurboCascadeNode(
			compressorEtaAd, piStag, etaT, lambdaOut,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
			etaM, precision,
		),
		burner: constructive.NewBurnerNode(fuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision),
	}

	result.linkPorts()

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

	return result
}

type gasGeneratorNode struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	burner       constructive.BurnerNode
	turboCascade TurboCascadeNode
}

func (node *gasGeneratorNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *gasGeneratorNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *gasGeneratorNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *gasGeneratorNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *gasGeneratorNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *gasGeneratorNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *gasGeneratorNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *gasGeneratorNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *gasGeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasGenerator")
}

func (node *gasGeneratorNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *gasGeneratorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
	}
}

func (node *gasGeneratorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *gasGeneratorNode) Burner() constructive.BurnerNode {
	return node.burner
}

func (node *gasGeneratorNode) TurboCascade() TurboCascadeNode {
	return node.turboCascade
}

func (node *gasGeneratorNode) Process() error {
	node.readInput()
	if err := node.turboCascade.ProcessCompressor(); err != nil {
		return err
	}
	if err := node.burner.Process(); err != nil {
		return err
	}
	if err := node.turboCascade.ProcessTurbine(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *gasGeneratorNode) readInput() {
	var cInput = node.turboCascade.CompressorComplexGasInput()
	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(), node.temperatureInput.GetState(),
			node.pressureInput.GetState(), node.massRateInput.GetState(),
		},
		[]graph.Port{
			cInput.GasInput(), cInput.TemperatureInput(), cInput.PressureInput(), cInput.MassRateInput(),
		},
	)
}

func (node *gasGeneratorNode) writeOutput() {
	var tOutput = node.turboCascade.TurbineComplexGasOutput()
	graph.SetAll(
		[]graph.PortState{
			tOutput.GasOutput().GetState(), tOutput.TemperatureOutput().GetState(),
			tOutput.PressureOutput().GetState(), tOutput.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
		},
	)
}

func (node *gasGeneratorNode) linkPorts() {
	nodes.LinkComplexOutToIn(node.turboCascade.CompressorComplexGasOutput(), node.burner)
	nodes.LinkComplexOutToIn(node.burner, node.turboCascade.TurbineComplexGasInput())
}
