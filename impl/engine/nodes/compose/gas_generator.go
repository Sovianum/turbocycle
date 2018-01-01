package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
)

type GasGeneratorNode interface {
	core.Node
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

	result.complexGasInput = core.NewAttachedPort(result)
	result.complexGasOutput = core.NewAttachedPort(result)

	return result
}

type gasGeneratorNode struct {
	core.BaseNode

	complexGasInput  core.Port
	complexGasOutput core.Port

	burner       constructive.BurnerNode
	turboCascade TurboCascadeNode
}

func (node *gasGeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasGenerator")
}

func (node *gasGeneratorNode) GetPorts() []core.Port {
	return []core.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *gasGeneratorNode) GetRequirePorts() []core.Port {
	return []core.Port{node.complexGasInput}
}

func (node *gasGeneratorNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.complexGasOutput}
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

func (node *gasGeneratorNode) ComplexGasInput() core.Port {
	return node.complexGasInput
}

func (node *gasGeneratorNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput
}

func (node *gasGeneratorNode) readInput() {
	node.turboCascade.CompressorComplexGasInput().SetState(node.complexGasInput.GetState())
}

func (node *gasGeneratorNode) writeOutput() {
	node.complexGasOutput.SetState(node.turboCascade.TurbineComplexGasOutput().GetState())
}

func (node *gasGeneratorNode) linkPorts() {
	core.Link(node.turboCascade.CompressorComplexGasOutput(), node.burner.ComplexGasInput())
	core.Link(node.burner.ComplexGasOutput(), node.turboCascade.TurbineComplexGasInput())
}
