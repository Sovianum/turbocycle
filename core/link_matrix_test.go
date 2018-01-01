package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraphMatrix_UnconnectedPorts(t *testing.T) {
	var nodes = []Node{NewTestNode(1, 0, true)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)
	assert.Equal(t, unconnectedPortsMsg, err.Error())
}

func TestNewGraphMatrix_ContextUndefined(t *testing.T) {
	var nodes = []Node{NewTestNode(0, 0, false)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)
	assert.Equal(t, contextUndefinedNodes, err.Error())
}

func TestGraphMatrix_getFreeNodes(t *testing.T) {
	var node1 = NewTestNode(0, 1, true)
	var node2 = NewTestNode(1, 0, true)

	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var freeNodes = matrix.getFreeNodes()
	assert.Equal(t, 1, len(freeNodes))
	assert.Equal(t, node1, freeNodes[0])
}

func TestGraphMatrix_GetCallOrder_OK(t *testing.T) {
	var node1 = NewTestNode(0, 1, true)
	var node2 = NewTestNode(1, 0, true)
	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var callOrder, orderErr = matrix.GetCallOrder()
	assert.Nil(t, orderErr)
	assert.Equal(t, 2, len(callOrder))
	assert.Equal(t, node1, callOrder[0])
	assert.Equal(t, node2, callOrder[1])
}

func TestGraphMatrix_GetCallOrder_Cyclic(t *testing.T) {
	var node1 = NewTestNode(1, 1, true)
	var node2 = NewTestNode(1, 1, true)
	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])
	Link(node2.GetUpdatePorts()[0], node1.GetRequirePorts()[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var _, orderErr = matrix.GetCallOrder()
	assert.NotNil(t, orderErr)
	assert.Equal(t, inaccessibleNodesMsg, orderErr.Error())
	assert.Equal(t, 2, len(orderErr.Nodes()))
}

func NewTestNode(requirePortNum, updatePortNum int, contextDefined bool) Node {
	var node = &testNode{
		requirePortNum: requirePortNum,
		updatePortNum:  updatePortNum,
		contextDefined: contextDefined,
		requirePorts:   make([]Port, requirePortNum),
		updatePorts:    make([]Port, updatePortNum),
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

type testNode struct {
	requirePortNum int
	updatePortNum  int

	contextDefined bool

	requirePorts []Port
	updatePorts  []Port
}

func (node *testNode) Process() error {
	panic("implement me")
}

func (node *testNode) GetRequirePorts() []Port {
	return node.requirePorts
}

func (node *testNode) GetUpdatePorts() []Port {
	return node.updatePorts
}

func (node *testNode) GetPorts() []Port {
	return append(node.requirePorts, node.updatePorts...)
}

func (node *testNode) GetPortByTag(tag string) (Port, error) {
	panic("implement me")
}

func (node *testNode) ContextDefined() bool {
	return node.contextDefined
}
