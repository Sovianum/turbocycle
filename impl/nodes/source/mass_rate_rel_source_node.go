package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type MassRateRelSourceNode interface {
	core.Node
	nodes.MassRateRelSource
}

type massRateRelSourceNode struct {
	ports core.PortsType
}

func NewMassRateRelSinkNode() MassRateRelSourceNode {
	var result = &massRateRelSourceNode{
		ports: make(core.PortsType),
	}

	result.ports[nodes.MassRateRelOutput] = core.NewPort()
	result.ports[nodes.MassRateRelOutput].SetInnerNode(result)
	result.ports[nodes.MassRateRelOutput].SetState(states.NewMassRateRelPortState(1))

	return result
}

func (node *massRateRelSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MassRateInputState core.PortState `json:"mass_rate_output_state"`
	}{
		MassRateInputState: node.MassRateRelOutput().GetState(),
	})
}

func (node *massRateRelSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *massRateRelSourceNode) Process() error {
	return nil
}

func (node *massRateRelSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *massRateRelSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.MassRateRelOutput}, nil
}

func (node *massRateRelSourceNode) GetPortTags() []string {
	return []string{nodes.MassRateRelOutput}
}

func (node *massRateRelSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.MassRateRelOutput:
		return node.ports[nodes.MassRateRelOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of massRateRelSourceNode can not be found", tag))
	}
}

func (node *massRateRelSourceNode) ContextDefined() bool {
	return true
}

func (node *massRateRelSourceNode) MassRateRelOut() float64 {
	return node.ports[nodes.MassRateRelOutput].GetState().(states.MassRateRelPortState).MassRateRel
}

func (node *massRateRelSourceNode) MassRateRelOutput() core.Port {
	return node.ports[nodes.MassRateRelOutput]
}
