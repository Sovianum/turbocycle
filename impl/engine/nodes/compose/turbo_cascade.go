package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
)

type TurboCascadeNode interface {
	graph.Node
	CompressorComplexGasInput() graph.Port
	CompressorComplexGasOutput() graph.Port
	TurbineComplexGasInput() graph.Port
	TurbineComplexGasOutput() graph.Port
	Compressor() constructive.CompressorNode
	Turbine() constructive.TurbineNode
	Transmission() constructive.TransmissionNode
	ProcessCompressor() error
	ProcessTurbine() error
}

func NewTurboCascadeNode(
	compressorEtaAd, piStag float64,
	etaT, lambdaOut float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) TurboCascadeNode {
	var result = &turboCascadeNode{
		compressor: constructive.NewCompressorNode(compressorEtaAd, piStag, precision),
		turbine: constructive.NewBlockedTurbineNode(
			etaT, lambdaOut, precision,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
		),
		transmission: constructive.NewTransmissionNode(etaM),
		powerSink:    sink.NewPowerSinkNode(),
	}

	result.linkPorts()

	result.compressorComplexGasInput = graph.NewAttachedPort(result)
	result.compressorComplexGasOutput = graph.NewAttachedPort(result)
	result.turbineComplexGasInput = graph.NewAttachedPort(result)
	result.turbineComplexGasOutput = graph.NewAttachedPort(result)

	return result
}

type turboCascadeNode struct {
	graph.BaseNode

	compressorComplexGasInput  graph.Port
	compressorComplexGasOutput graph.Port
	turbineComplexGasInput     graph.Port
	turbineComplexGasOutput    graph.Port

	compressor   constructive.CompressorNode
	turbine      constructive.BlockedTurbineNode
	transmission constructive.TransmissionNode
	powerSink    nodes.PowerSink
}

func (node *turboCascadeNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TurboCascade")
}

func (node *turboCascadeNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.compressorComplexGasInput,
		node.compressorComplexGasOutput,
		node.turbineComplexGasInput,
		node.turbineComplexGasOutput,
	}
}

func (node *turboCascadeNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.compressorComplexGasInput,
		node.turbineComplexGasInput,
	}
}

func (node *turboCascadeNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.compressorComplexGasOutput,
		node.turbineComplexGasOutput,
	}
}

func (node *turboCascadeNode) Compressor() constructive.CompressorNode {
	return node.compressor
}

func (node *turboCascadeNode) Turbine() constructive.TurbineNode {
	return node.turbine
}

func (node *turboCascadeNode) Transmission() constructive.TransmissionNode {
	return node.transmission
}

func (node *turboCascadeNode) ProcessCompressor() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	node.writeOutput()

	return nil
}

func (node *turboCascadeNode) ProcessTurbine() error {
	node.readInput()
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) Process() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) TurbineComplexGasOutput() graph.Port {
	return node.turbineComplexGasOutput
}

func (node *turboCascadeNode) TurbineComplexGasInput() graph.Port {
	return node.turbineComplexGasInput
}

func (node *turboCascadeNode) CompressorComplexGasOutput() graph.Port {
	return node.compressorComplexGasOutput
}

func (node *turboCascadeNode) CompressorComplexGasInput() graph.Port {
	return node.compressorComplexGasInput
}

func (node *turboCascadeNode) linkPorts() {
	graph.Link(node.compressor.PowerOutput(), node.transmission.PowerInput())
	graph.Link(node.transmission.PowerOutput(), node.turbine.PowerInput())
	graph.Link(node.turbine.PowerOutput(), node.powerSink.PowerInput())
}

func (node *turboCascadeNode) readInput() {
	node.compressor.ComplexGasInput().SetState(node.compressorComplexGasInput.GetState())
	node.turbine.ComplexGasInput().SetState(node.turbineComplexGasInput.GetState())
}

func (node *turboCascadeNode) writeOutput() {
	node.compressorComplexGasOutput.SetState(node.compressor.ComplexGasOutput().GetState())
	node.turbineComplexGasOutput.SetState(node.turbine.ComplexGasOutput().GetState())
}
