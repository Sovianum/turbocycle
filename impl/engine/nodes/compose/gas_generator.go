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

	result.complexGasInput = graph.NewAttachedPort(result)
	result.complexGasOutput = graph.NewAttachedPort(result)

	return result
}

type gasGeneratorNode struct {
	graph.BaseNode

	complexGasInput  graph.Port
	complexGasOutput graph.Port

	burner       constructive.BurnerNode
	turboCascade TurboCascadeNode
}

func (node *gasGeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasGenerator")
}

func (node *gasGeneratorNode) GetPorts() []graph.Port {
	return []graph.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *gasGeneratorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{node.complexGasInput}
}

func (node *gasGeneratorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{node.complexGasOutput}
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

func (node *gasGeneratorNode) ComplexGasInput() graph.Port {
	return node.complexGasInput
}

func (node *gasGeneratorNode) ComplexGasOutput() graph.Port {
	return node.complexGasOutput
}

func (node *gasGeneratorNode) readInput() {
	node.turboCascade.CompressorComplexGasInput().SetState(node.complexGasInput.GetState())
}

func (node *gasGeneratorNode) writeOutput() {
	node.complexGasOutput.SetState(node.turboCascade.TurbineComplexGasOutput().GetState())
}

func (node *gasGeneratorNode) linkPorts() {
	graph.Link(node.turboCascade.CompressorComplexGasOutput(), node.burner.ComplexGasInput())
	graph.Link(node.burner.ComplexGasOutput(), node.turboCascade.TurbineComplexGasInput())
}
