package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type MassRateRelSourceNode interface {
	core.Node
	nodes.MassRateRelSource
}

type massRateRelSourceNode struct {
	ports       core.PortsType
	massRateRel float64
}

func NewMassRateRelSourceNode(massRateRel float64) MassRateRelSourceNode {
	var result = &massRateRelSourceNode{
		ports:       make(core.PortsType),
		massRateRel: massRateRel,
	}

	result.ports[nodes.MassRateRelOutput] = core.NewPort()
	result.ports[nodes.MassRateRelOutput].SetInnerNode(result)
	result.ports[nodes.MassRateRelOutput].SetState(states.NewMassRateRelPortState(massRateRel))

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
	node.ports[nodes.MassRateRelOutput].SetState(states.NewMassRateRelPortState(node.massRateRel))
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
		return nil, fmt.Errorf("Port %s of massRateRelSourceNode can not be found", tag)
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
