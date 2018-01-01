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
)

func NewThreeShaftsCoolingRegeneratorScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	middlePressureCascade compose.TurboCascadeNode,
	cooler constructive.CoolerNode,
	regenerativeGasGenerator compose.RegenerativeGasGeneratorNode,
	middlePressureCompressorPipe constructive.PressureLossNode,
	highPressureTurbinePipe constructive.PressureLossNode,
	middlePressureTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,
) ThreeShaftsCoolingRegeneratorScheme {
	return &threeShaftsCoolingRegeneratorScheme{
		gasSource:                    gasSource,
		inletPressureDrop:            inletPressureDrop,
		middlePressureCascade:        middlePressureCascade,
		regenerativeGasGenerator:     regenerativeGasGenerator,
		middlePressureCompressorPipe: middlePressureCompressorPipe,
		cooler:                    cooler,
		highPressureTurbinePipe:   highPressureTurbinePipe,
		middlePressureTurbinePipe: middlePressureTurbinePipe,
		freeTurbineBlock:          freeTurbineBlock,
		gasSink:                   sink.NewComplexGasSinkNode(),
		powerSink:                 sink.NewPowerSinkNode(),
		breaker1:                  helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
		breaker2:                  helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
	}
}

type ThreeShaftsCoolingRegeneratorScheme interface {
	Scheme
	DoubleCompressor
	InitGasGeneratorCompressor(state states.ComplexGasPortState)
	InitGasGeneratorHeatExchanger(state states.ComplexGasPortState)
}

type threeShaftsCoolingRegeneratorScheme struct {
	gasSource                    source.ComplexGasSourceNode
	inletPressureDrop            constructive.PressureLossNode
	middlePressureCascade        compose.TurboCascadeNode
	cooler                       constructive.CoolerNode
	regenerativeGasGenerator     compose.RegenerativeGasGeneratorNode
	middlePressureCompressorPipe constructive.PressureLossNode
	highPressureTurbinePipe      constructive.PressureLossNode
	middlePressureTurbinePipe    constructive.PressureLossNode
	freeTurbineBlock             compose.FreeTurbineBlockNode
	gasSink                      sink.ComplexGasSinkNode
	powerSink                    nodes.PowerSink
	breaker1                     helper.CycleBreakNode
	breaker2                     helper.CycleBreakNode
}

func (scheme *threeShaftsCoolingRegeneratorScheme) LowPressureCompressor() constructive.CompressorNode {
	return scheme.middlePressureCascade.Compressor()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) HighPressureCompressor() constructive.CompressorNode {
	return scheme.regenerativeGasGenerator.TurboCascade().Compressor()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) InitGasGeneratorCompressor(state states.ComplexGasPortState) {
	scheme.breaker1.DataSourcePort().SetState(state)
}

func (scheme *threeShaftsCoolingRegeneratorScheme) InitGasGeneratorHeatExchanger(state states.ComplexGasPortState) {
	scheme.breaker2.DataSourcePort().SetState(state)
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.regenerativeGasGenerator.Burner().ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.regenerativeGasGenerator.Burner().GetFuelRateRel() * massRateRel
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetQLower() float64 {
	return scheme.regenerativeGasGenerator.Burner().Fuel().QLower()
}

func (scheme *threeShaftsCoolingRegeneratorScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.middlePressureCascade, scheme.cooler,
		scheme.middlePressureCompressorPipe, scheme.regenerativeGasGenerator, scheme.highPressureTurbinePipe,
		scheme.middlePressureTurbinePipe, scheme.freeTurbineBlock, scheme.gasSink, scheme.breaker1, scheme.breaker2,
	})
}

func (scheme *threeShaftsCoolingRegeneratorScheme) linkPorts() {
	graph.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	graph.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.middlePressureCascade.CompressorComplexGasInput())

	graph.Link(scheme.middlePressureCascade.CompressorComplexGasOutput(), scheme.middlePressureCompressorPipe.ComplexGasInput())
	graph.Link(scheme.middlePressureCompressorPipe.ComplexGasOutput(), scheme.cooler.ComplexGasInput())
	graph.Link(scheme.cooler.ComplexGasOutput(), scheme.breaker1.DataSourcePort())
	graph.Link(scheme.breaker1.UpdatePort(), scheme.regenerativeGasGenerator.ComplexGasInput())
	graph.Link(scheme.regenerativeGasGenerator.ComplexGasOutput(), scheme.highPressureTurbinePipe.ComplexGasInput())
	graph.Link(scheme.highPressureTurbinePipe.ComplexGasOutput(), scheme.middlePressureCascade.TurbineComplexGasInput())
	graph.Link(scheme.middlePressureCascade.TurbineComplexGasOutput(), scheme.middlePressureTurbinePipe.ComplexGasInput())
	graph.Link(scheme.middlePressureTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	graph.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
	graph.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.breaker2.DataSourcePort())
	graph.Link(scheme.breaker2.UpdatePort(), scheme.regenerativeGasGenerator.HeatExchangerHotInput())
	graph.Link(scheme.regenerativeGasGenerator.HeatExchangerHotOutput(), scheme.gasSink.ComplexGasInput())
}
