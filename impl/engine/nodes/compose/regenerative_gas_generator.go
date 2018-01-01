package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
)

type RegenerativeGasGeneratorNode interface {
	core.Node
	nodes.ComplexGasChannel
	HeatExchangerHotInput() core.Port
	HeatExchangerHotOutput() core.Port
	Burner() constructive.BurnerNode
	TurboCascade() TurboCascadeNode
	Regenerator() constructive.RegeneratorNode
}

func NewRegenerativeGasGeneratorNode(
	compressorEtaAd, piStag float64,
	fuel fuel.GasFuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0 float64,
	etaT, lambdaOut float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(constructive.TurbineNode) float64,
	sigmaRegenerator float64,
	sigmaRegeneratorPipe float64,
	etaM float64,
	precision float64,
) RegenerativeGasGeneratorNode {
	var result = &regenerativeGasGeneratorNode{
		turboCascade: NewTurboCascadeNode(
			compressorEtaAd, piStag, etaT, lambdaOut,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
			etaM, precision,
		),
		burner:                  constructive.NewBurnerNode(fuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision),
		regenerator:             constructive.NewRegeneratorNode(sigmaRegenerator, precision, constructive.SigmaByColdSide),
		regeneratorPressureDrop: constructive.NewPressureLossNode(sigmaRegeneratorPipe),
	}

	result.complexGasInput = core.NewAttachedPort(result)
	result.complexGasOutput = core.NewAttachedPort(result)
	result.heatExchangerHotInput = core.NewAttachedPort(result)
	result.heatExchangerHotOutput = core.NewAttachedPort(result)

	result.linkPorts()
	return result
}

type regenerativeGasGeneratorNode struct {
	core.BaseNode

	complexGasInput        core.Port
	complexGasOutput       core.Port
	heatExchangerHotInput  core.Port
	heatExchangerHotOutput core.Port

	burner                  constructive.BurnerNode
	turboCascade            TurboCascadeNode
	regenerator             constructive.RegeneratorNode
	regeneratorPressureDrop constructive.PressureLossNode
}

func (node *regenerativeGasGeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "RegenerativeGasGenerator")
}

func (node *regenerativeGasGeneratorNode) GetPorts() []core.Port {
	return []core.Port{node.complexGasInput, node.complexGasOutput, node.heatExchangerHotInput, node.heatExchangerHotOutput}
}

func (node *regenerativeGasGeneratorNode) GetRequirePorts() []core.Port {
	return []core.Port{node.complexGasInput, node.heatExchangerHotInput}
}

func (node *regenerativeGasGeneratorNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.complexGasOutput, node.heatExchangerHotOutput}
}

func (node *regenerativeGasGeneratorNode) Burner() constructive.BurnerNode {
	return node.burner
}

func (node *regenerativeGasGeneratorNode) TurboCascade() TurboCascadeNode {
	return node.turboCascade
}

func (node *regenerativeGasGeneratorNode) Regenerator() constructive.RegeneratorNode {
	return node.regenerator
}

func (node *regenerativeGasGeneratorNode) Process() error {
	node.readInput()
	if err := node.turboCascade.ProcessCompressor(); err != nil {
		return err
	}
	if err := node.regeneratorPressureDrop.Process(); err != nil {
		return err
	}
	if err := node.regenerator.Process(); err != nil {
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

func (node *regenerativeGasGeneratorNode) ContextDefined() bool {
	return true
}

func (node *regenerativeGasGeneratorNode) ComplexGasInput() core.Port {
	return node.complexGasInput
}

func (node *regenerativeGasGeneratorNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotInput() core.Port {
	return node.heatExchangerHotInput
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotOutput() core.Port {
	return node.heatExchangerHotOutput
}

func (node *regenerativeGasGeneratorNode) readInput() {
	node.turboCascade.CompressorComplexGasInput().SetState(node.complexGasInput.GetState())
	node.regenerator.HotInput().SetState(node.heatExchangerHotInput.GetState())
}

func (node *regenerativeGasGeneratorNode) writeOutput() {
	node.complexGasOutput.SetState(node.turboCascade.TurbineComplexGasOutput().GetState())
	node.heatExchangerHotOutput.SetState(node.regenerator.HotOutput().GetState())
}

func (node *regenerativeGasGeneratorNode) linkPorts() {
	core.Link(node.turboCascade.CompressorComplexGasOutput(), node.regeneratorPressureDrop.ComplexGasInput())
	core.Link(node.regeneratorPressureDrop.ComplexGasOutput(), node.regenerator.ColdInput())
	core.Link(node.regenerator.ColdOutput(), node.burner.ComplexGasInput())
	core.Link(node.burner.ComplexGasOutput(), node.turboCascade.TurbineComplexGasInput())
}
