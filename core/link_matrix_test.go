package core

import (
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
)

func TestNewGraphMatrix_UnconnectedPorts(t *testing.T) {
	var nodes = []Node{newTestNodeWithoutAction(1, 0, true)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)
	assert.Equal(t, unconnectedPortsMsg, err.Error())
}

func TestNewGraphMatrix_ContextUndefined(t *testing.T) {
	var nodes = []Node{newTestNodeWithoutAction(0, 0, false)}
	var _, err = newGraphMatrix(nodes)

	assert.NotNil(t, err)
	assert.Equal(t, contextUndefinedNodes, err.Error())
}

func TestGraphMatrix_getFreeNodes(t *testing.T) {
	var node1 = newTestNodeWithoutAction(0, 1, true)
	var node2 = newTestNodeWithoutAction(1, 0, true)

	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var freeNodes = matrix.getFreeNodes()
	assert.Equal(t, 1, len(freeNodes))
	assert.Equal(t, node1, freeNodes[0])
}

func TestGraphMatrix_GetCallOrder_OK(t *testing.T) {
	var node1 = newTestNodeWithoutAction(0, 1, true)
	var node2 = newTestNodeWithoutAction(1, 0, true)
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
	var node1 = newTestNodeWithoutAction(1, 1, true)
	var node2 = newTestNodeWithoutAction(1, 1, true)
	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])
	Link(node2.GetUpdatePorts()[0], node1.GetRequirePorts()[0])

	var matrix, err = newGraphMatrix([]Node{node1, node2})
	assert.Nil(t, err)

	var _, orderErr = matrix.GetCallOrder()
	assert.NotNil(t, orderErr)
	assert.Equal(t, inaccessibleNodesMsg, orderErr.Error())
	assert.Equal(t, 2, len(orderErr.Nodes()))
}

func newTestNodeWithoutAction(requirePortNum, updatePortNum int, contextDefined bool) *testNode {
	return newTestNode(requirePortNum, updatePortNum, contextDefined, func() error {
		return nil
	})
}

func newTestNode(requirePortNum, updatePortNum int, contextDefined bool, action func() error) *testNode {
	var node = &testNode{
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

type testNode struct {
	BaseNode

	requirePortNum int
	updatePortNum  int

	contextDefined bool

	action func() error

	requirePorts []Port
	updatePorts  []Port
}

func (node *testNode) GetName() string {
	return common.EitherString(node.name, "test_node")
}

func (node *testNode) Process() error {
	return node.action()
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

func newTestPortState(val float64) testPortState {
	return testPortState{val: val}
}

type testPortState struct {
	val float64
}

func (state testPortState) Mix(another PortState, relaxCoef float64) (PortState, error) {
	return testPortState{
		val: common.Lerp(state.val, another.(testPortState).val, relaxCoef),
	}, nil
}

func (state testPortState) MaxResidual(another PortState) (float64, error) {
	return common.GetRelResidual(state.val, another.(testPortState).val), nil
}
