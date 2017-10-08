package schemes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
)

func NewGtn16TwoShafts(
	gasSource source.ComplexGasSourceNode,
	inletPressureDrop constructive.PressureLossNode,
	gasGenerator compose.GasGeneratorNode,
	compressorTurbinePipe constructive.PressureLossNode,
	freeTurbineBlock compose.FreeTurbineBlockNode,

) Gtn16TwoShafts {
	return &gtn16TwoShafts{
		gasSource:             gasSource,
		inletPressureDrop:     inletPressureDrop,
		gasGenerator:          gasGenerator,
		compressorTurbinePipe: compressorTurbinePipe,
		freeTurbineBlock:      freeTurbineBlock,
		gasSink:               sink.NewComplexGasSinkNode(),
		powerSink:             sink.NewPowerSinkNode(),
	}
}

type Gtn16TwoShafts interface {
	Scheme
}

type gtn16TwoShafts struct {
	gasSource             source.ComplexGasSourceNode
	inletPressureDrop     constructive.PressureLossNode
	gasGenerator          compose.GasGeneratorNode
	compressorTurbinePipe constructive.PressureLossNode
	freeTurbineBlock      compose.FreeTurbineBlockNode
	gasSink               sink.ComplexGasSinkNode
	powerSink             nodes.PowerSink
}

func (scheme *gtn16TwoShafts) GetNetwork() core.Network {
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
