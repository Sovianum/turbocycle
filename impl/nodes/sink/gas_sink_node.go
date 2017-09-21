package sink

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type GasSinkNode interface {
	core.Node
	nodes.GasSink
}

type gasSinkNode struct {
	ports core.PortsType
}

func NewGasSinkNode() GasSinkNode {
	var result = &gasSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.GasInput] = core.NewPort()
	result.ports[nodes.GasInput].SetInnerNode(result)
	result.ports[nodes.GasInput].SetState(states.NewGasPortState(gases.GetAir()))

	return result
}

func (node *gasSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState core.PortState `json:"gas_input_state"`
	}{
		GasInputState: node.GasInput().GetState(),
	})
}

func (node *gasSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasSinkNode) Process() error {
	return nil
}

func (node *gasSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.GasInput}, nil
}

func (node *gasSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *gasSinkNode) GetPortTags() []string {
	return []string{nodes.GasInput}
}

func (node *gasSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.GasInput:
		return node.ports[nodes.GasInput], nil
	default:
		return nil, fmt.Errorf("Port %s of gasSinkNode can not be found", tag)
	}
}

func (node *gasSinkNode) ContextDefined() bool {
	return true
}

func (node *gasSinkNode) GasInput() core.Port {
	return node.ports[nodes.GasInput]
}
