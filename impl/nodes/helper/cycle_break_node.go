package helper

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
)

type CycleBreakNode interface {
	core.Node
	UpdatePort() core.Port
	DataSourcePort() core.Port
}

type initializerNode struct {
	ports core.PortsType
}

func NewCycleBreakerNode(initialState core.PortState) CycleBreakNode {
	var result = &initializerNode{ports: make(core.PortsType)}

	result.ports[nodes.UpdatePort] = core.NewPort()
	result.ports[nodes.UpdatePort].SetInnerNode(result)
	result.ports[nodes.UpdatePort].SetState(initialState)

	result.ports[nodes.DataSourcePort] = core.NewPort()
	result.ports[nodes.DataSourcePort].SetInnerNode(result)
	result.ports[nodes.DataSourcePort].SetState(initialState)

	return result
}

func (node *initializerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		State core.PortState
	}{
		State: node.DataSourcePort().GetState(),
	})
}

func (node *initializerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *initializerNode) Process() error {
	node.UpdatePort().SetState(node.DataSourcePort().GetState())
	return nil
}

func (node *initializerNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *initializerNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.UpdatePort}, nil
}

func (node *initializerNode) GetPortTags() []string {
	return []string{nodes.UpdatePort, nodes.DataSourcePort}
}

func (node *initializerNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.UpdatePort:
		return node.UpdatePort(), nil
	case nodes.DataSourcePort:
		return node.DataSourcePort(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found in cycle breaker", tag)
	}
}

func (node *initializerNode) ContextDefined() bool {
	return true
}

func (node *initializerNode) UpdatePort() core.Port {
	return node.ports[nodes.UpdatePort]
}

func (node *initializerNode) DataSourcePort() core.Port {
	return node.ports[nodes.DataSourcePort]
}
