package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func NewTwoShaftsScheme(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	gasGenerator compose.GasGeneratorNode,
	compressorTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,

) TwoShaftsScheme {
	return &twoShaftsScheme{
		gasSource:             gasSource,
		inletPressureDrop:     inletPressureDrop,
		gasGenerator:          gasGenerator,
		compressorTurbinePipe: compressorTurbinePipe,
		freeTurbineBlock:      freeTurbineBlock,
		gasSink:               sink.NewComplexGasSinkNode(),
		powerSink:             sink.NewPowerSinkNode(),
	}
}

type TwoShaftsScheme interface {
	Scheme
	SingleCompressor
}

type twoShaftsScheme struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	gasGenerator          compose.GasGeneratorNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode
	gasSink               sink.ComplexGasSinkNode
	powerSink             nodes.PowerSink
}

func (scheme *twoShaftsScheme) Compressor() constructive.CompressorNode {
	return scheme.gasGenerator.TurboCascade().Compressor()
}

func (scheme *twoShaftsScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return lSpecific * massRateRel
}

func (scheme *twoShaftsScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.gasGenerator.Burner().ComplexGasInput().GetState().(states.ComplexGasPortState).MassRateRel
	return scheme.gasGenerator.Burner().GetFuelRateRel() * massRateRel
}

func (scheme *twoShaftsScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *twoShaftsScheme) GetNetwork() (core.Network, core.GraphError) {
	scheme.linkPorts()

	return core.NewNetwork([]core.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.gasGenerator,
		scheme.compressorTurbinePipe, scheme.freeTurbineBlock, scheme.gasSink,
	})
}

func (scheme *twoShaftsScheme) linkPorts() {
	core.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	core.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.gasGenerator.ComplexGasInput())
	core.Link(scheme.gasGenerator.ComplexGasOutput(), scheme.compressorTurbinePipe.ComplexGasInput())
	core.Link(scheme.compressorTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.gasSink.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())
}
