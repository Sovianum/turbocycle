package schemes

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewTwoShaftsRegeneratorScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	turboCascade compose.TurboCascadeNode,
	burner constructive.BurnerNode,
	compressorTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
	regenerator constructive.RegeneratorNode,

) TwoShaftsRegeneratorScheme {
	return &twoShaftsRegeneratorScheme{
		gasSource:             gasSource,
		inletPressureDrop:     inletPressureDrop,
		turboCascade:          turboCascade,
		burner:                burner,
		compressorTurbinePipe: compressorTurbinePipe,
		freeTurbineBlock:      freeTurbineBlock,
		regenerator:           regenerator,
		breaker1: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
		breaker2: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
	}
}

type TwoShaftsRegeneratorScheme interface {
	TwoShaftsScheme
	Regenerator() constructive.RegeneratorNode
}

type twoShaftsRegeneratorScheme struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	turboCascade          compose.TurboCascadeNode
	burner                constructive.BurnerNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode
	regenerator           constructive.RegeneratorNode

	gasSink         sink.SinkNode
	temperatureSink sink.SinkNode
	pressureSink    sink.SinkNode
	massRateSink    sink.SinkNode
	powerSink       sink.SinkNode

	breaker1 helper.ComplexCycleBreakNode
	breaker2 helper.ComplexCycleBreakNode
}

func (scheme *twoShaftsRegeneratorScheme) Regenerator() constructive.RegeneratorNode {
	return scheme.regenerator
}

func (scheme *twoShaftsRegeneratorScheme) FreeTurbineBlock() compose.FreeTurbineBlockNode {
	return scheme.freeTurbineBlock
}

func (scheme *twoShaftsRegeneratorScheme) CompressorTurbinePipe() constructive.PressureLossNode {
	return scheme.compressorTurbinePipe
}

func (scheme *twoShaftsRegeneratorScheme) Burner() constructive.BurnerNode {
	return scheme.burner
}

func (scheme *twoShaftsRegeneratorScheme) TurboCascade() compose.TurboCascadeNode {
	return scheme.turboCascade
}

func (scheme *twoShaftsRegeneratorScheme) InletPressureDrop() constructive.PressureLossNode {
	return scheme.inletPressureDrop
}

func (scheme *twoShaftsRegeneratorScheme) GasSource() source.ComplexGasSourceNode {
	return scheme.gasSource
}

func (scheme *twoShaftsRegeneratorScheme) Compressor() constructive.CompressorNode {
	return scheme.turboCascade.Compressor()
}

func (scheme *twoShaftsRegeneratorScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return lSpecific * massRateRel
}

func (scheme *twoShaftsRegeneratorScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.burner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return scheme.burner.FuelRateRel() * massRateRel
}

func (scheme *twoShaftsRegeneratorScheme) GetQLower() float64 {
	return scheme.burner.Fuel().QLower()
}

func (scheme *twoShaftsRegeneratorScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.turboCascade,
		scheme.regenerator, scheme.burner, scheme.compressorTurbinePipe,
		scheme.freeTurbineBlock,
		scheme.gasSink, scheme.temperatureSink, scheme.pressureSink, scheme.massRateSink, scheme.powerSink,
		scheme.breaker1, scheme.breaker2,
	})
}

func (scheme *twoShaftsRegeneratorScheme) linkPorts() {
	nodes.LinkComplexOutToIn(scheme.gasSource, scheme.inletPressureDrop)
	nodes.LinkComplexOutToIn(scheme.inletPressureDrop, scheme.turboCascade.CompressorComplexGasInput())
	nodes.LinkComplexOutToIn(scheme.turboCascade.CompressorComplexGasOutput(), scheme.regenerator.ColdInput())
	nodes.LinkComplexOutToIn(scheme.regenerator.ColdOutput(), scheme.breaker1)
	nodes.LinkComplexOutToIn(scheme.breaker1, scheme.burner)
	nodes.LinkComplexOutToIn(scheme.burner, scheme.turboCascade.TurbineComplexGasInput())
	nodes.LinkComplexOutToIn(scheme.turboCascade.TurbineComplexGasOutput(), scheme.compressorTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.compressorTurbinePipe, scheme.freeTurbineBlock)
	nodes.LinkComplexOutToIn(scheme.freeTurbineBlock, scheme.breaker2)
	nodes.LinkComplexOutToIn(scheme.breaker2, scheme.regenerator.HotInput())

	var hotOutput = scheme.regenerator.HotOutput()
	scheme.gasSink = sink.SinkPort(hotOutput.GasOutput())
	scheme.temperatureSink = sink.SinkPort(hotOutput.TemperatureOutput())
	scheme.pressureSink = sink.SinkPort(hotOutput.PressureOutput())
	scheme.massRateSink = sink.SinkPort(hotOutput.MassRateOutput())
	scheme.powerSink = sink.SinkPort(scheme.freeTurbineBlock.PowerOutput())
}
