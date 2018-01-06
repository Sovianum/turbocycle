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

func NewThreeShaftsScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	gasGenerator compose.GasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsScheme {
	return &threeShaftsScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		gasGenerator:                 gasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		highPressureTurbinePipe:      highPressureTurbinePipe,
		middlePressureTurbinePipe:    middlePressureTurbinePipe,
		freeTurbineBlock:             freeTurbineBlock,
		breaker: helper.NewComplexCycleBreakNode(
			gases.GetAir(), states.StandardTemperature, states.StandardPressure, 1,
		),
	}
}

type ThreeShaftsScheme interface {
	Scheme
	DoubleCompressor
	MainBurner() constructive.BurnerNode
	HighPressureTurbine() constructive.TurbineNode
	LowPressureTurbine() constructive.TurbineNode
	FreeTurbine() constructive.TurbineNode
	FreeTurbineBlock() compose.FreeTurbineBlockNode
	MiddlePressureTurbinePipe() constructive.PressureLossNode
	HighPressureTurbinePipe() constructive.PressureLossNode
	MiddlePressureCompressorPipe() constructive.PressureLossNode
	GasGenerator() compose.GasGeneratorNode
	MiddlePressureCascade() compose.TurboCascadeNode
	InletPressureDrop() constructive.PressureLossNode
	GasSource() source.ComplexGasSourceNode
}

type threeShaftsScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	gasGenerator                 compose.GasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode

	gasSink         sink.SinkNode
	pressureSink    sink.SinkNode
	temperatureSink sink.SinkNode
	massRateSink    sink.SinkNode
	powerSink       sink.SinkNode

	breaker helper.ComplexCycleBreakNode
}

func (scheme *threeShaftsScheme) MainBurner() constructive.BurnerNode {
	return scheme.gasGenerator.Burner()
}

func (scheme *threeShaftsScheme) HighPressureTurbine() constructive.TurbineNode {
	return scheme.gasGenerator.TurboCascade().Turbine()
}

func (scheme *threeShaftsScheme) LowPressureTurbine() constructive.TurbineNode {
	return scheme.middlePressureCascade.Turbine()
}

func (scheme *threeShaftsScheme) FreeTurbine() constructive.TurbineNode {
	return scheme.freeTurbineBlock.FreeTurbine()
}

func (scheme *threeShaftsScheme) FreeTurbineBlock() compose.FreeTurbineBlockNode {
	return scheme.freeTurbineBlock
}

func (scheme *threeShaftsScheme) MiddlePressureTurbinePipe() constructive.PressureLossNode {
	return scheme.middlePressureTurbinePipe
}

func (scheme *threeShaftsScheme) HighPressureTurbinePipe() constructive.PressureLossNode {
	return scheme.highPressureTurbinePipe
}

func (scheme *threeShaftsScheme) MiddlePressureCompressorPipe() constructive.PressureLossNode {
	return scheme.middlePressureCompressorPipe
}

func (scheme *threeShaftsScheme) GasGenerator() compose.GasGeneratorNode {
	return scheme.gasGenerator
}

func (scheme *threeShaftsScheme) MiddlePressureCascade() compose.TurboCascadeNode {
	return scheme.middlePressureCascade
}

func (scheme *threeShaftsScheme) InletPressureDrop() constructive.PressureLossNode {
	return scheme.inletPressureDrop
}

func (scheme *threeShaftsScheme) GasSource() source.ComplexGasSourceNode {
	return scheme.gasSource
}

func (scheme *threeShaftsScheme) LowPressureCompressor() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsScheme) HighPressureCompressor() constructive.CompressorNode {
	return scheme.gasGenerator.TurboCascade().Compressor()
}

func (scheme *threeShaftsScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return lSpecific * massRateRel
}

func (scheme *threeShaftsScheme) GetFuelMassRateRel() float64 {
	var burner = scheme.gasGenerator.Burner()
	var massRateRel = burner.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return burner.FuelRateRel() * massRateRel
}

func (scheme *threeShaftsScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.middlePressureCascade,
		scheme.middlePressureCompressorPipe, scheme.gasGenerator, scheme.highPressureTurbinePipe,
		scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock,
		scheme.gasSink, scheme.temperatureSink, scheme.pressureSink, scheme.massRateSink, scheme.powerSink,
		scheme.breaker,
	})
}

func (scheme *threeShaftsScheme) linkPorts() {
	nodes.LinkComplexOutToIn(scheme.gasSource, scheme.inletPressureDrop)
	nodes.LinkComplexOutToIn(scheme.inletPressureDrop, scheme.middlePressureCascade.CompressorComplexGasInput())

	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureCompressorPipe, scheme.breaker)
	nodes.LinkComplexOutToIn(scheme.breaker, scheme.gasGenerator)
	nodes.LinkComplexOutToIn(scheme.gasGenerator, scheme.highPressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.highPressureTurbinePipe, scheme.middlePressureCascade.TurbineComplexGasInput())
	nodes.LinkComplexOutToIn(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock)

	scheme.powerSink = sink.SinkPort(scheme.freeTurbineBlock.PowerOutput())
	scheme.gasSink = sink.SinkPort(scheme.freeTurbineBlock.GasOutput())
	scheme.pressureSink = sink.SinkPort(scheme.freeTurbineBlock.PressureOutput())
	scheme.temperatureSink = sink.SinkPort(scheme.freeTurbineBlock.TemperatureOutput())
	scheme.massRateSink = sink.SinkPort(scheme.freeTurbineBlock.MassRateOutput())
}
