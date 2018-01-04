package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/material/fuel"
)

type RegenerativeGasGeneratorNode interface {
	graph.Node
	CompressorInput() nodes.ComplexGasSink
	TurbineOutput() nodes.ComplexGasSource
	HeatExchangerHotInput() nodes.ComplexGasSink
	HeatExchangerHotOutput() nodes.ComplexGasSource
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

	graph.AttachAllPorts(
		result,
		&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
		&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,

		&result.regeneratorHotTemperatureInput, &result.regeneratorHotPressureInput,
		&result.regeneratorHotGasInput, &result.regeneratorHotMassRateInput,

		&result.regeneratorHotTemperatureOutput, &result.regeneratorHotPressureOutput,
		&result.regeneratorHotGasOutput, &result.regeneratorHotMassRateOutput,
	)

	result.linkPorts()
	return result
}

type regenerativeGasGeneratorNode struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	regeneratorHotTemperatureInput graph.Port
	regeneratorHotPressureInput    graph.Port
	regeneratorHotGasInput         graph.Port
	regeneratorHotMassRateInput    graph.Port

	regeneratorHotTemperatureOutput graph.Port
	regeneratorHotPressureOutput    graph.Port
	regeneratorHotGasOutput         graph.Port
	regeneratorHotMassRateOutput    graph.Port

	burner                  constructive.BurnerNode
	turboCascade            TurboCascadeNode
	regenerator             constructive.RegeneratorNode
	regeneratorPressureDrop constructive.PressureLossNode
}

func (node *regenerativeGasGeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "RegenerativeGasGenerator")
}

func (node *regenerativeGasGeneratorNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,

		node.regeneratorHotTemperatureInput, node.regeneratorHotPressureInput,
		node.regeneratorHotGasInput, node.regeneratorHotMassRateInput,

		node.regeneratorHotTemperatureOutput, node.regeneratorHotPressureOutput,
		node.regeneratorHotGasOutput, node.regeneratorHotMassRateOutput,
	}
}

func (node *regenerativeGasGeneratorNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,

		node.regeneratorHotTemperatureInput, node.regeneratorHotPressureInput,
		node.regeneratorHotGasInput, node.regeneratorHotMassRateInput,
	}
}

func (node *regenerativeGasGeneratorNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,

		node.regeneratorHotTemperatureOutput, node.regeneratorHotPressureOutput,
		node.regeneratorHotGasOutput, node.regeneratorHotMassRateOutput,
	}
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

func (node *regenerativeGasGeneratorNode) CompressorInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput,
	)
}

func (node *regenerativeGasGeneratorNode) TurbineOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
	)
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotInput() nodes.ComplexGasSink {
	return helper.NewPseudoComplexGasSink(
		node.regeneratorHotGasInput,
		node.regeneratorHotTemperatureInput,
		node.regeneratorHotPressureInput,
		node.regeneratorHotMassRateInput,
	)
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotOutput() nodes.ComplexGasSource {
	return helper.NewPseudoComplexGasSource(
		node.regeneratorHotGasOutput,
		node.regeneratorHotTemperatureOutput,
		node.regeneratorHotPressureOutput,
		node.regeneratorHotMassRateOutput,
	)
}

func (node *regenerativeGasGeneratorNode) readInput() {
	var c = node.turboCascade.CompressorComplexGasInput()
	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(), node.temperatureInput.GetState(),
			node.pressureInput.GetState(), node.massRateInput.GetState(),
		},
		[]graph.Port{
			c.GasInput(), c.TemperatureInput(), c.PressureInput(), c.MassRateInput(),
		},
	)

	var rSink = node.regenerator.HotInput()
	graph.SetAll(
		[]graph.PortState{
			node.regeneratorHotGasInput.GetState(), node.regeneratorHotTemperatureInput.GetState(),
			node.regeneratorHotPressureInput.GetState(), node.regeneratorHotMassRateInput.GetState(),
		},
		[]graph.Port{
			rSink.GasInput(), rSink.TemperatureInput(), rSink.PressureInput(), rSink.MassRateInput(),
		},
	)
}

func (node *regenerativeGasGeneratorNode) writeOutput() {
	var tSource = node.turboCascade.TurbineComplexGasOutput()
	graph.SetAll(
		[]graph.PortState{
			tSource.GasOutput().GetState(), tSource.TemperatureOutput().GetState(),
			tSource.PressureOutput().GetState(), tSource.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
		},
	)

	var hSource = node.regenerator.HotOutput()
	graph.SetAll(
		[]graph.PortState{
			hSource.GasOutput().GetState(), hSource.TemperatureOutput().GetState(),
			hSource.PressureOutput().GetState(), hSource.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.regeneratorHotGasOutput, node.regeneratorHotTemperatureOutput,
			node.regeneratorHotPressureOutput, node.regeneratorHotMassRateOutput,
		},
	)
}

func (node *regenerativeGasGeneratorNode) linkPorts() {
	nodes.LinkComplexOutToIn(node.turboCascade.CompressorComplexGasOutput(), node.regeneratorPressureDrop)
	nodes.LinkComplexOutToIn(node.regeneratorPressureDrop, node.regenerator.ColdInput())
	nodes.LinkComplexOutToIn(node.regenerator.ColdOutput(), node.burner)
	nodes.LinkComplexOutToIn(node.burner, node.turboCascade.TurbineComplexGasInput())
}
