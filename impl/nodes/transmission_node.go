package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type TransmissionNode interface {
	core.Node
	PowerChannel
}

type transmissionNode struct {
	ports core.PortsType
	etaM  float64
}

func NewTransmissionNode(etaM float64) TransmissionNode {
	var transmissionNode = &transmissionNode{
		ports: make(core.PortsType),
		etaM:  etaM,
	}

	var inputPort = core.NewPort()
	inputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[powerInput] = inputPort
	transmissionNode.ports[powerInput].SetState(states.StandartPowerState())

	var outputPort = core.NewPort()
	outputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[powerOutput] = outputPort
	transmissionNode.ports[powerOutput].SetState(states.StandartPowerState())

	return transmissionNode
}

func (node *transmissionNode) ContextDefined() bool {
	return true
}

func (node *transmissionNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case powerInput:
		return node.powerInput(), nil
	case powerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *transmissionNode) GetRequirePortTags() ([]string, error) {
	return []string{powerInput}, nil
}

func (node *transmissionNode) GetUpdatePortTags() ([]string, error) {
	return []string{powerOutput}, nil
}

func (node *transmissionNode) GetPortTags() []string {
	return []string{powerInput, powerOutput}
}

func (node *transmissionNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *transmissionNode) PowerInput() core.Port {
	return node.powerInput()
}

func (node *transmissionNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *transmissionNode) Process() error {
	var inputState = node.PowerInput().GetState()

	switch v := inputState.(type) {
	case states.PowerPortState:
		var state = inputState.(states.PowerPortState)
		state.LSpecific /= node.etaM // TODO check if division is a correct operation (depends on ifo flow direction)
		node.powerOutput().SetState(state)
		return nil
	default:
		return common.GetTypeError("PowerPortState", v)
	}
}

func (node *transmissionNode) powerInput() core.Port {
	return node.ports[powerInput]
}

func (node *transmissionNode) powerOutput() core.Port {
	return node.ports[powerOutput]
}
