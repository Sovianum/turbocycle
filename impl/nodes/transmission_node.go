package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/common"
)

type TransmissionNode struct {
	ports core.PortsType
	etaM float64
}

func NewTransmissionNode(etaM float64) *TransmissionNode {
	var transmissionNode = &TransmissionNode{
		ports:make(map[string]*core.Port),
		etaM:etaM,
	}

	var inputPort = core.NewPort()
	inputPort.SetDest(transmissionNode)
	transmissionNode.ports[powerInput] = inputPort

	var outputPort = core.NewPort()
	outputPort.SetSrc(transmissionNode)
	transmissionNode.ports[powerOutput] = outputPort

	return transmissionNode
}

func (node *TransmissionNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *TransmissionNode) PowerInput() *core.Port {
	return node.ports[powerInput]
}

func (node *TransmissionNode) PowerOutput() *core.Port {
	return node.ports[powerOutput]
}

func (node *TransmissionNode) Process() error {
	var inputState = node.PowerInput().GetState()

	switch v := inputState.(type) {
	case states.PowerPortState:
		var state = inputState.(states.PowerPortState)
		state.LSpecific /= node.etaM	// TODO check if division is a correct operation (depends on ifo flow direction)
		node.PowerOutput().SetState(state)
		return nil
	default:
		return common.GetTypeError("PowerPortState", v)
	}
}





