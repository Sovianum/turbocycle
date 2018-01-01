package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
)

type FreeTurbineBlockNode interface {
	core.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
	FreeTurbine() constructive.FreeTurbineNode
	OutletPressureLoss() constructive.PressureLossNode
}

func NewFreeTurbineBlock(
	pAtm float64,
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(node constructive.TurbineNode) float64,
	pressureLossSigma float64,
) FreeTurbineBlockNode {
	var result = &freeTurbineBlockNode{
		atmNode: source.NewComplexGasSourceNode(nil, 0, pAtm), // first two arguments are not used cos they will be sent to sinks
		turbine: constructive.NewFreeTurbineNode(
			etaT, lambdaOut, precision,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
		),
		pressureLoss: constructive.NewPressureLossNode(pressureLossSigma),
		assembler:    helper.NewGasStateAssemblerNode(),
		disassembler: helper.NewGasStateDisassemblerNode(),
		tSink:        sink.NewTemperatureSinkNode(),
		gSink:        sink.NewGasSinkNode(),
		mSink:        sink.NewMassRateRelSinkNode(),
		hub:          helper.NewHubNode(),
	}
	result.linkPorts()

	result.complexGasInput = core.NewAttachedPort(result)
	result.complexGasOutput = core.NewAttachedPort(result)
	result.powerOutput = core.NewAttachedPort(result)

	return result
}

type freeTurbineBlockNode struct {
	core.BaseNode

	complexGasInput  core.Port
	complexGasOutput core.Port
	powerOutput      core.Port

	atmNode      source.ComplexGasSourceNode
	turbine      constructive.FreeTurbineNode
	pressureLoss constructive.PressureLossNode
	assembler    helper.GasStateAssemblerNode
	disassembler helper.GasStateAssemblerNode
	tSink        nodes.TemperatureSink
	gSink        nodes.GasSink
	mSink        nodes.MassRateRelSink
	hub          helper.HubNode
}

func (node *freeTurbineBlockNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "FreeTurbineBlock")
}

func (node *freeTurbineBlockNode) GetPorts() []core.Port {
	return []core.Port{node.complexGasInput, node.complexGasOutput, node.powerOutput}
}

func (node *freeTurbineBlockNode) GetRequirePorts() []core.Port {
	return []core.Port{node.complexGasInput}
}

func (node *freeTurbineBlockNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.complexGasOutput, node.powerOutput}
}

func (node *freeTurbineBlockNode) FreeTurbine() constructive.FreeTurbineNode {
	return node.turbine
}

func (node *freeTurbineBlockNode) OutletPressureLoss() constructive.PressureLossNode {
	return node.pressureLoss
}

func (node *freeTurbineBlockNode) Process() error {
	node.readInput()
	if err := node.atmNode.Process(); err != nil {
		return err
	}
	if err := node.pressureLoss.Process(); err != nil {
		return err
	}
	if err := node.disassembler.Process(); err != nil {
		return err
	}
	if err := node.hub.Process(); err != nil {
		return err
	}
	if err := node.turbine.Process(); err != nil {
		return err
	}
	if err := node.assembler.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *freeTurbineBlockNode) ComplexGasInput() core.Port {
	return node.complexGasInput
}

func (node *freeTurbineBlockNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput
}

func (node *freeTurbineBlockNode) PowerOutput() core.Port {
	return node.powerOutput
}

func (node *freeTurbineBlockNode) linkPorts() {
	core.Link(node.atmNode.ComplexGasOutput(), node.pressureLoss.ComplexGasOutput())
	core.Link(node.pressureLoss.ComplexGasInput(), node.disassembler.ComplexGasPort())

	core.Link(node.disassembler.TemperaturePort(), node.tSink.TemperatureInput())
	core.Link(node.disassembler.MassRateRelPort(), node.mSink.MassRateRelInput())
	core.Link(node.disassembler.GasPort(), node.gSink.GasInput())
	core.Link(node.disassembler.PressurePort(), node.hub.Inlet())

	core.Link(node.hub.Outlet1(), node.turbine.PressureOutput())
	core.Link(node.hub.Outlet2(), node.assembler.PressurePort())

	core.Link(node.turbine.TemperatureOutput(), node.assembler.TemperaturePort())
	core.Link(node.turbine.GasOutput(), node.assembler.GasPort())
	core.Link(node.turbine.MassRateRelOutput(), node.assembler.MassRateRelPort())
}

func (node *freeTurbineBlockNode) readInput() {
	node.turbine.ComplexGasInput().SetState(node.complexGasInput.GetState())
}

func (node *freeTurbineBlockNode) writeOutput() {
	node.complexGasOutput.SetState(node.assembler.ComplexGasPort().GetState())
	node.powerOutput.SetState(node.turbine.PowerOutput().GetState())
}
