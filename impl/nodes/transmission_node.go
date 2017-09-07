package nodes

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type transmissionNode struct {
	ports core.PortsType
	etaM  float64
}

func NewTransmissionNode(etaM float64) *transmissionNode {
	var transmissionNode = &transmissionNode{
		ports: make(map[string]*core.Port),
		etaM:  etaM,
	}

	var inputPort = core.NewPort()
	inputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[powerInput] = inputPort

	var outputPort = core.NewPort()
	outputPort.SetInnerNode(transmissionNode)
	transmissionNode.ports[powerOutput] = outputPort

	return transmissionNode
}

func (node *transmissionNode) GetRequiredPorts() []string {
	return []string{powerInput}
}

func (node *transmissionNode) GetUpdatedPorts() []string {
	return []string{powerOutput}
}

func (node *transmissionNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *transmissionNode) PowerInput() *core.Port {
	return node.powerInput()
}

func (node *transmissionNode) PowerOutput() *core.Port {
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

func (node *transmissionNode) powerInput() *core.Port {
	return node.ports[powerInput]
}

func (node *transmissionNode) powerOutput() *core.Port {
	return node.ports[powerOutput]
}
