package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"encoding/json"
	"errors"
	"fmt"
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

	result.ports[portA] = core.NewPort()
	result.ports[portA].SetInnerNode(result)
	result.ports[portA].SetState(initialState)

	result.ports[portB] = core.NewPort()
	result.ports[portB].SetInnerNode(result)
	result.ports[portB].SetState(initialState)

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
	return []string{portA, portB}, nil
}

func (node *initializerNode) GetPortTags() []string {
	return []string{portA, portB}
}

func (node *initializerNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case portA:
		return node.PortA(), nil
	case portB:
		return node.PortB(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found in cycle breaker", tag))
	}
}

func (node *initializerNode) ContextDefined() bool {
	return true
}

func (node *initializerNode) PortA() core.Port {
	return node.ports[portA]
}

func (node *initializerNode) PortB() core.Port {
	return node.ports[portB]
}
