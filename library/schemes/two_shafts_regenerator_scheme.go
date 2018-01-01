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
		gasSink:               sink.NewComplexGasSinkNode(),
		powerSink:             sink.NewPowerSinkNode(),
		breaker1:              helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
		breaker2:              helper.NewCycleBreakerNode(states.StandardAtmosphereState()),
	}
}

type TwoShaftsRegeneratorScheme interface {
	Scheme
	SingleCompressor
}

type twoShaftsRegeneratorScheme struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	turboCascade          compose.TurboCascadeNode
	burner                constructive.BurnerNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode
	regenerator           constructive.RegeneratorNode
	gasSink               sink.ComplexGasSinkNode
	powerSink             nodes.PowerSink
	breaker1              helper.CycleBreakNode
	breaker2              helper.CycleBreakNode
}

func (scheme *twoShaftsRegeneratorScheme) Compressor() constructive.CompressorNode {
	return scheme.turboCascade.Compressor()
}

func (scheme *twoShaftsRegeneratorScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *twoShaftsRegeneratorScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.burner.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.burner.GetFuelRateRel() * massRateRel
}

func (scheme *twoShaftsRegeneratorScheme) GetQLower() float64 {
	return scheme.burner.Fuel().QLower()
}

func (scheme *twoShaftsRegeneratorScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.turboCascade,
		scheme.regenerator, scheme.burner, scheme.compressorTurbinePipe,
		scheme.freeTurbineBlock, scheme.gasSink, scheme.breaker1, scheme.breaker2,
	})
}

func (scheme *twoShaftsRegeneratorScheme) linkPorts() {
	graph.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	graph.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.turboCascade.CompressorComplexGasInput())
	graph.Link(scheme.turboCascade.CompressorComplexGasOutput(), scheme.regenerator.ColdInput())
	graph.Link(scheme.regenerator.ColdOutput(), scheme.breaker1.DataSourcePort())
	graph.Link(scheme.breaker1.UpdatePort(), scheme.burner.ComplexGasInput())
	graph.Link(scheme.burner.ComplexGasOutput(), scheme.turboCascade.TurbineComplexGasInput())
	graph.Link(scheme.turboCascade.TurbineComplexGasOutput(), scheme.compressorTurbinePipe.ComplexGasInput())
	graph.Link(scheme.compressorTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	graph.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.breaker2.DataSourcePort())
	graph.Link(scheme.breaker2.UpdatePort(), scheme.regenerator.HotInput())
	graph.Link(scheme.regenerator.HotOutput(), scheme.gasSink.ComplexGasInput())
	graph.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
}
