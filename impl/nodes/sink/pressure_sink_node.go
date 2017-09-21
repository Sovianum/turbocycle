package sink

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PressureSinkNode interface {
	core.Node
	nodes.PressureSink
}

type pressureSinkNode struct {
	ports core.PortsType
}

func NewPressureSinkNode() PressureSinkNode {
	var result = &pressureSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.PressureInput] = core.NewPort()
	result.ports[nodes.PressureInput].SetInnerNode(result)
	result.ports[nodes.PressureInput].SetState(states.NewPressurePortState(1e5))

	return result
}

func (node *pressureSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PressureInputState core.PortState `json:"pressure_input_state"`
	}{
		PressureInputState: node.PressureInput().GetState(),
	})
}

func (node *pressureSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *pressureSinkNode) Process() error {
	return nil
}

func (node *pressureSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.PressureInput}, nil
}

func (node *pressureSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *pressureSinkNode) GetPortTags() []string {
	return []string{nodes.PressureInput}
}

func (node *pressureSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PressureInput:
		return node.ports[nodes.PressureInput], nil
	default:
		return nil, fmt.Errorf("Port %s of pressureSinkNode can not be found", tag)
	}
}

func (node *pressureSinkNode) ContextDefined() bool {
	return true
}

func (node *pressureSinkNode) PStagIn() float64 {
	return node.ports[nodes.PressureInput].GetState().(states.PressurePortState).PStag
}

func (node *pressureSinkNode) PressureInput() core.Port {
	return node.ports[nodes.PressureInput]
}
