package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PowerSourceNode interface {
	core.Node
	nodes.PowerSource
}

type powerSourceNode struct {
	ports     core.PortsType
	lSpecific float64
}

func NewPowerSourceNode(lRel float64) PowerSourceNode {
	var result = &powerSourceNode{
		ports:     make(core.PortsType),
		lSpecific: lRel,
	}

	result.ports[nodes.PowerOutput] = core.NewPort()
	result.ports[nodes.PowerOutput].SetInnerNode(result)
	result.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	return result
}

func (node *powerSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PowerOutputState core.PortState `json:"power_source_state"`
	}{
		PowerOutputState: node.PowerOutput().GetState(),
	})
}

func (node *powerSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *powerSourceNode) Process() error {
	node.ports[nodes.PowerOutput].SetState(states.NewPowerPortState(node.lSpecific))
	return nil
}

func (node *powerSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *powerSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.PowerOutput}, nil
}

func (node *powerSourceNode) GetPortTags() []string {
	return []string{nodes.PowerOutput}
}

func (node *powerSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PowerOutput:
		return node.ports[nodes.PowerOutput], nil
	default:
		return nil, fmt.Errorf("Port %s of powerSourceNode can not be found", tag)
	}
}

func (node *powerSourceNode) ContextDefined() bool {
	return true
}

func (node *powerSourceNode) PowerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
