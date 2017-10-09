package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
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

func (scheme *twoShaftsScheme) GetNetwork() core.Network {
	var nodeMap = make(map[string]core.Node)
	nodeMap[inputGasSourceName] = scheme.gasSource
	nodeMap[inletPressureDropName] = scheme.inletPressureDrop
	nodeMap[gasGeneratorName] = scheme.gasGenerator
	nodeMap[compressorTurbinePipeName] = scheme.compressorTurbinePipe
	nodeMap[freeTurbineBlockName] = scheme.freeTurbineBlock
	nodeMap[outputGasSinkName] = scheme.gasSink

	core.Link(scheme.gasSource.ComplexGasOutput(), scheme.inletPressureDrop.ComplexGasInput())
	core.Link(scheme.inletPressureDrop.ComplexGasOutput(), scheme.gasGenerator.ComplexGasInput())
	core.Link(scheme.gasGenerator.ComplexGasOutput(), scheme.compressorTurbinePipe.ComplexGasInput())
	core.Link(scheme.compressorTurbinePipe.ComplexGasOutput(), scheme.freeTurbineBlock.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.ComplexGasOutput(), scheme.gasSink.ComplexGasInput())
	core.Link(scheme.freeTurbineBlock.PowerOutput(), scheme.powerSink.PowerInput())

	return core.NewNetwork(nodeMap)
}
