package compose

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func NewFreeTurbineBlock(
	pAtm float64,
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(node constructive.TurbineNode) float64,
	pressureLossSigma float64,
) FreeTurbineBlockNode {
	var result = &freeTurbineBlockNode{
		ports:   make(core.PortsType),
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
		hub:          helper.NewHubNode(states.StandardAtmosphereState()),
	}
	result.linkPorts()

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.PowerOutput] = core.NewPort()
	result.ports[nodes.PowerOutput].SetInnerNode(result)
	result.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	return result
}

type FreeTurbineBlockNode interface {
	core.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
	FreeTurbine() constructive.FreeTurbineNode
	OutletPressureLoss() constructive.PressureLossNode
}

type freeTurbineBlockNode struct {
	ports        core.PortsType
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

func (node *freeTurbineBlockNode) FreeTurbine() constructive.FreeTurbineNode {
	return node.turbine
}

func (node *freeTurbineBlockNode) OutletPressureLoss() constructive.PressureLossNode {
	return node.pressureLoss
}

func (node *freeTurbineBlockNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState               `json:"gas_input_state"`
		GasOutputState   core.PortState               `json:"gas_output_state"`
		PowerOutputState core.PortState               `json:"power_output_state"`
		Turbine          constructive.FreeTurbineNode `json:"turbine"`
	}{
		GasInputState:    node.complexGasInput().GetState(),
		GasOutputState:   node.complexGasOutput().GetState(),
		PowerOutputState: node.powerOutput().GetState(),
		Turbine:          node.turbine,
	})
}

func (node *freeTurbineBlockNode) GetPorts() core.PortsType {
	return node.ports
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

func (node *freeTurbineBlockNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *freeTurbineBlockNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput, nodes.PowerOutput}, nil
}

func (node *freeTurbineBlockNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.ComplexGasOutput, nodes.PowerOutput}
}

func (node *freeTurbineBlockNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.ComplexGasOutput:
		return node.complexGasOutput(), nil
	case nodes.PowerOutput:
		return node.powerOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *freeTurbineBlockNode) ContextDefined() bool {
	return true
}

func (node *freeTurbineBlockNode) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *freeTurbineBlockNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput()
}

func (node *freeTurbineBlockNode) PowerOutput() core.Port {
	return node.powerOutput()
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
	node.turbine.ComplexGasInput().SetState(node.ComplexGasInput().GetState())
}

func (node *freeTurbineBlockNode) writeOutput() {
	node.complexGasOutput().SetState(node.assembler.ComplexGasPort().GetState())
	node.PowerOutput().SetState(node.turbine.PowerOutput().GetState())
}

func (node *freeTurbineBlockNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *freeTurbineBlockNode) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *freeTurbineBlockNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
