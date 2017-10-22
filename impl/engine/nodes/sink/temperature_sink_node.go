package sink

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type TemperatureSinkNode interface {
	core.Node
	nodes.TemperatureSink
}

type temperatureSinkNode struct {
	ports core.PortsType
}

func NewTemperatureSinkNode() TemperatureSinkNode {
	var result = &temperatureSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.TemperatureInput] = core.NewPort()
	result.ports[nodes.TemperatureInput].SetInnerNode(result)
	result.ports[nodes.TemperatureInput].SetState(states.NewTemperaturePortState(288))

	return result
}

func (node *temperatureSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TemperatureInputState core.PortState `json:"temperature_input_state"`
	}{
		TemperatureInputState: node.TemperatureInput().GetState(),
	})
}

func (node *temperatureSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *temperatureSinkNode) Process() error {
	return nil
}

func (node *temperatureSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.TemperatureInput}, nil
}

func (node *temperatureSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *temperatureSinkNode) GetPortTags() []string {
	return []string{nodes.TemperatureInput}
}

func (node *temperatureSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.TemperatureInput:
		return node.ports[nodes.TemperatureInput], nil
	default:
		return nil, fmt.Errorf("Port %s of temperatureSinkNode can not be found", tag)
	}
}

func (node *temperatureSinkNode) ContextDefined() bool {
	return true
}

func (node *temperatureSinkNode) TStagIn() float64 {
	return node.ports[nodes.TemperatureInput].GetState().(states.TemperaturePortState).TStag
}

func (node *temperatureSinkNode) TemperatureInput() core.Port {
	return node.ports[nodes.TemperatureInput]
}
