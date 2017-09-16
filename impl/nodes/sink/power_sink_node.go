package sink

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PowerSinkNode interface {
	core.Node
	nodes.PowerSink
}

type powerSinkNode struct {
	ports core.PortsType
}

func NewPowerSinkNode() PowerSinkNode {
	var result = &powerSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.PowerInput] = core.NewPort()
	result.ports[nodes.PowerInput].SetInnerNode(result)
	result.ports[nodes.PowerInput].SetState(states.StandardPowerState())

	return result
}

func (node *powerSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PowerInputState core.PortState `json:"power_input_state"`
	}{
		PowerInputState: node.PowerInput().GetState(),
	})
}

func (node *powerSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *powerSinkNode) Process() error {
	return nil
}

func (node *powerSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.PowerInput}, nil
}

func (node *powerSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSinkNode) GetPortTags() []string {
	return []string{nodes.PowerInput}
}

func (node *powerSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PowerInput:
		return node.ports[nodes.PowerInput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of powerSinkNode can not be found", tag))
	}
}

func (node *powerSinkNode) ContextDefined() bool {
	return true
}

func (node *powerSinkNode) PowerInput() core.Port {
	return node.ports[nodes.PowerInput]
}
