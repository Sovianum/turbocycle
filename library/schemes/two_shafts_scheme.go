package schemes

import (
	"github.com/Sovianum/turbocycle/core/graph"
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
	}
}

type TwoShaftsScheme interface {
	Scheme
	SingleCompressor
	FreeTurbineBlock() compose.FreeTurbineBlockNode
	CompressorTurbinePipe() constructive.PressureLossNode
	Burner() constructive.BurnerNode
	TurboCascade() compose.TurboCascadeNode
	InletPressureDrop() constructive.PressureLossNode
	GasSource() source.ComplexGasSourceNode
}

type twoShaftsScheme struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	gasGenerator          compose.GasGeneratorNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode

	gasSink         sink.SinkNode
	temperatureSink sink.SinkNode
	pressureSink    sink.SinkNode
	massRateSink    sink.SinkNode
	powerSink       sink.SinkNode
}

func (scheme *twoShaftsScheme) Burner() constructive.BurnerNode {
	return scheme.gasGenerator.Burner()
}

func (scheme *twoShaftsScheme) TurboCascade() compose.TurboCascadeNode {
	return scheme.gasGenerator.TurboCascade()
}

func (scheme *twoShaftsScheme) FreeTurbineBlock() compose.FreeTurbineBlockNode {
	return scheme.freeTurbineBlock
}

func (scheme *twoShaftsScheme) CompressorTurbinePipe() constructive.PressureLossNode {
	return scheme.compressorTurbinePipe
}

func (scheme *twoShaftsScheme) InletPressureDrop() constructive.PressureLossNode {
	return scheme.inletPressureDrop
}

func (scheme *twoShaftsScheme) GasSource() source.ComplexGasSourceNode {
	return scheme.gasSource
}

func (scheme *twoShaftsScheme) Compressor() constructive.CompressorNode {
	return scheme.gasGenerator.TurboCascade().Compressor()
}

func (scheme *twoShaftsScheme) GetSpecificPower() float64 {
	var turbine = scheme.freeTurbineBlock.FreeTurbine()
	var lSpecific = turbine.PowerOutput().GetState().(states.PowerPortState).LSpecific
	var massRateRel = turbine.MassRateInput().GetState().(states.MassRatePortState).MassRate
	return lSpecific * massRateRel
}

func (scheme *twoShaftsScheme) GetFuelMassRateRel() float64 {
	var massRateRel = scheme.gasGenerator.Burner().MassRateInput().GetState().(states.MassRatePortState).MassRate
	return scheme.gasGenerator.Burner().FuelRateRel() * massRateRel
}

func (scheme *twoShaftsScheme) GetQLower() float64 {
	return scheme.gasGenerator.Burner().Fuel().QLower()
}

func (scheme *twoShaftsScheme) GetNetwork() (graph.Network, graph.GraphError) {
	scheme.linkPorts()

	return graph.NewNetwork([]graph.Node{
		scheme.gasSource, scheme.inletPressureDrop, scheme.gasGenerator,
		scheme.compressorTurbinePipe, scheme.freeTurbineBlock,
		scheme.gasSink, scheme.temperatureSink, scheme.pressureSink, scheme.massRateSink, scheme.powerSink,
	})
}

func (scheme *twoShaftsScheme) linkPorts() {
	nodes.LinkComplexOutToIn(scheme.gasSource, scheme.inletPressureDrop)
	nodes.LinkComplexOutToIn(scheme.inletPressureDrop, scheme.gasGenerator)
	nodes.LinkComplexOutToIn(scheme.gasGenerator, scheme.compressorTurbinePipe)
	nodes.LinkComplexOutToIn(scheme.compressorTurbinePipe, scheme.freeTurbineBlock)

	var turbine = scheme.freeTurbineBlock
	scheme.gasSink = sink.SinkPort(turbine.GasOutput())
	scheme.temperatureSink = sink.SinkPort(turbine.TemperatureOutput())
	scheme.pressureSink = sink.SinkPort(turbine.PressureOutput())
	scheme.massRateSink = sink.SinkPort(turbine.MassRateOutput())
	scheme.powerSink = sink.SinkPort(scheme.freeTurbineBlock.PowerOutput())
}
