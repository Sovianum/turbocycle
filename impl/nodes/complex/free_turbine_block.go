package complex

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/nodes/helper"
	"github.com/Sovianum/turbocycle/impl/nodes/source"
)

type FreeTurbineBlockNode interface {
	core.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
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
	return nil
}

func (node *freeTurbineBlockNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.GasInput}, nil
}

func (node *freeTurbineBlockNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput, nodes.PowerOutput}, nil
}

func (node *freeTurbineBlockNode) GetPortTags() []string {
	return []string{nodes.GasInput, nodes.ComplexGasOutput, nodes.PowerOutput}
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

func (node *freeTurbineBlockNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *freeTurbineBlockNode) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *freeTurbineBlockNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
