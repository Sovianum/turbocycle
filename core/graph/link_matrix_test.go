package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraphMatrix_UnconnectedPorts(t *testing.T) {
	var nodes = []Node{NewTestNodeWithoutAction(1, 0, true)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)

	var ports, _ = nodes[0].GetRequirePorts()
	assert.Equal(t, getUnconnectedErrMsg(ports), err.Error())
}

func TestNewGraphMatrix_ContextUndefined(t *testing.T) {
	var nodes = []Node{NewTestNodeWithoutAction(0, 0, false)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)
	assert.Equal(t, getUndefinedErrMsg(nodes), err.Error())
}

func TestGraphMatrix_getFreeNodes(t *testing.T) {
	var node1 = NewTestNodeWithoutAction(0, 1, true)
	var node2 = NewTestNodeWithoutAction(1, 0, true)

	Link(node1.updatePorts[0], node2.requirePorts[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var freeNodes = matrix.getFreeNodes()
	assert.Equal(t, 1, len(freeNodes))
	assert.Equal(t, node1, freeNodes[0])
}

func TestGraphMatrix_GetCallOrder_OK(t *testing.T) {
	var node1 = NewTestNodeWithoutAction(0, 1, true)
	var node2 = NewTestNodeWithoutAction(1, 0, true)
	Link(node1.updatePorts[0], node2.requirePorts[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var callOrder, orderErr = matrix.GetCallOrder()
	assert.Nil(t, orderErr)
	assert.Equal(t, 2, len(callOrder))
	assert.Equal(t, node1, callOrder[0])
	assert.Equal(t, node2, callOrder[1])
}

func TestGraphMatrix_GetCallOrder_Cyclic(t *testing.T) {
	var node1 = NewTestNodeWithoutAction(1, 1, true)
	var node2 = NewTestNodeWithoutAction(1, 1, true)
	Link(node1.updatePorts[0], node2.requirePorts[0])
	Link(node2.updatePorts[0], node1.requirePorts[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var _, orderErr = matrix.GetCallOrder()
	assert.NotNil(t, orderErr)
	assert.Equal(t, getInaccessibleErrMsg([]Node{node1, node2}), orderErr.Error())
	assert.Equal(t, 2, len(orderErr.Nodes()))
}
