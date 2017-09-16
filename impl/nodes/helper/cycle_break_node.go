package helper

import (
	"github.com/Sovianum/turbocycle/core"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/impl/nodes"
)

type CycleBreakNode interface {
	core.Node
	PortA() core.Port
	PortB() core.Port
}

type initializerNode struct {
	ports core.PortsType
}

func NewCycleBreakerNode(initialState core.PortState) CycleBreakNode {
	var result = &initializerNode{ports:make(core.PortsType)}

	result.ports[nodes.PortA] = core.NewPort()
	result.ports[nodes.PortA].SetInnerNode(result)
	result.ports[nodes.PortA].SetState(initialState)

	result.ports[nodes.PortB] = core.NewPort()
	result.ports[nodes.PortB].SetInnerNode(result)
	result.ports[nodes.PortB].SetState(initialState)

	return result
}

func (node *initializerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PortAState core.PortState
		PortBState core.PortState
	}{
		PortAState:node.PortA().GetState(),
		PortBState:node.PortB().GetState(),
	})
}

func (node *initializerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *initializerNode) Process() error {
	var oldStateA = node.PortA().GetState()
	var oldStateB = node.PortB().GetState()

	node.PortA().SetState(oldStateB)
	node.PortB().SetState(oldStateA)

	return nil
}

func (node *initializerNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *initializerNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.PortA, nodes.PortB}, nil
}

func (node *initializerNode) GetPortTags() []string {
	return []string{nodes.PortA, nodes.PortB}
}

func (node *initializerNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.PortA:
		return node.PortA(), nil
	case nodes.PortB:
		return node.PortB(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found in cycle breaker", tag))
	}
}

func (node *initializerNode) ContextDefined() bool {
	return true
}

func (node *initializerNode) PortA() core.Port {
	return node.ports[nodes.PortA]
}

func (node *initializerNode) PortB() core.Port {
	return node.ports[nodes.PortB]
}
