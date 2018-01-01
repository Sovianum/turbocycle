package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork_Solve_OK(t *testing.T) {
	var node1 = newTestNodeWithoutAction(0, 1, true)
	node1.action = func() error {
		node1.GetUpdatePorts()[0].SetState(newTestPortState(1.))
		return nil
	}

	var node2 = newTestNodeWithoutAction(1, 1, true)
	Link(node1.GetUpdatePorts()[0], node2.GetRequirePorts()[0])
	node2.action = func() error {
		var state = newTestPortState(
			node2.GetRequirePorts()[0].GetState().(testPortState).val + 1,
		)
		node2.GetUpdatePorts()[0].SetState(state)
		return nil
	}

	var node3 = newTestNodeWithoutAction(1, 0, true)
	Link(node2.GetUpdatePorts()[0], node3.GetRequirePorts()[0])

	var network, networkErr = NewNetwork([]Node{node1, node2, node3})
	assert.Nil(t, networkErr)

	var converged, err = network.Solve(1, 1, 100, 1e-7)
	assert.True(t, converged)
	assert.Nil(t, err)

	var state = node3.GetRequirePorts()[0].GetState().(testPortState)
	assert.InDelta(t, 2., state.val, 1e-9)
}
