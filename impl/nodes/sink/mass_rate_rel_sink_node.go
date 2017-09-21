package sink

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type MassRateRelSinkNode interface {
	core.Node
	nodes.MassRateRelSink
}

type massRateRelSinkNode struct {
	ports core.PortsType
}

func NewMassRateRelSinkNode() MassRateRelSinkNode {
	var result = &massRateRelSinkNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.MassRateRelInput] = core.NewPort()
	result.ports[nodes.MassRateRelInput].SetInnerNode(result)
	result.ports[nodes.MassRateRelInput].SetState(states.NewMassRateRelPortState(1))

	return result
}

func (node *massRateRelSinkNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MassRateInputState core.PortState `json:"mass_rate_input_state"`
	}{
		MassRateInputState: node.MassRateRelInput().GetState(),
	})
}

func (node *massRateRelSinkNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *massRateRelSinkNode) Process() error {
	return nil
}

func (node *massRateRelSinkNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.MassRateRelInput}, nil
}

func (node *massRateRelSinkNode) GetUpdatePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *massRateRelSinkNode) GetPortTags() []string {
	return []string{nodes.MassRateRelInput}
}

func (node *massRateRelSinkNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.MassRateRelInput:
		return node.ports[nodes.MassRateRelInput], nil
	default:
		return nil, fmt.Errorf("Port %s of massRateRelSinkNode can not be found", tag)
	}
}

func (node *massRateRelSinkNode) ContextDefined() bool {
	return true
}

func (node *massRateRelSinkNode) MassRateRelIn() float64 {
	return node.ports[nodes.MassRateRelInput].GetState().(states.MassRateRelPortState).MassRateRel
}

func (node *massRateRelSinkNode) MassRateRelInput() core.Port {
	return node.ports[nodes.MassRateRelInput]
}
