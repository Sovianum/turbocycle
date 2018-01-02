package graph

import (
	"github.com/Sovianum/turbocycle/common"
)

func NewTestNodeWithoutAction(requirePortNum, updatePortNum int, contextDefined bool) *TestNode {
	return NewTestNode(requirePortNum, updatePortNum, contextDefined, func() error {
		return nil
	})
}

func NewTestNode(requirePortNum, updatePortNum int, contextDefined bool, action func() error) *TestNode {
	var node = &TestNode{
		requirePortNum: requirePortNum,
		updatePortNum:  updatePortNum,
		contextDefined: contextDefined,
		requirePorts:   make([]Port, requirePortNum),
		updatePorts:    make([]Port, updatePortNum),
		action:         action,
	}

	for i := 0; i != requirePortNum; i++ {
		node.requirePorts[i] = NewAttachedPort(node)
	}

	for i := 0; i != updatePortNum; i++ {
		var port = NewPort()
		port.SetInnerNode(node)
		node.updatePorts[i] = NewAttachedPort(node)
	}

	return node
}

type TestNode struct {
	BaseNode

	requirePortNum int
	updatePortNum  int

	contextDefined bool

	action func() error

	requirePorts []Port
	updatePorts  []Port
}

func (node *TestNode) SetAction(action func() error) {
	node.action = action
}

func (node *TestNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "test_node")
}

func (node *TestNode) Process() error {
	return node.action()
}

func (node *TestNode) GetRequirePorts() []Port {
	return node.requirePorts
}

func (node *TestNode) GetUpdatePorts() []Port {
	return node.updatePorts
}

func (node *TestNode) GetPorts() []Port {
	return append(node.requirePorts, node.updatePorts...)
}

func (node *TestNode) GetPortByTag(tag string) (Port, error) {
	panic("implement me")
}

func (node *TestNode) ContextDefined() bool {
	return node.contextDefined
}
