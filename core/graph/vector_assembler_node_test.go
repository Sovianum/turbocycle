package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorAssemblerNode_GetPairPort(t *testing.T) {
	var node = NewVectorAssemblerNode()
	var port1 = NewPort()
	var port2 = NewPort()

	var pairPort1 = node.GetPairPort(port1)
	var pairPort2 = node.GetPairPort(port2)

	var requirePorts = node.GetRequirePorts()

	assert.Equal(t, 2, len(requirePorts))
	assert.Equal(t, pairPort1, requirePorts[0])
	assert.Equal(t, pairPort2, requirePorts[1])
}

func TestVectorAssemblerNode_DeletePairPort(t *testing.T) {
	var node = NewVectorAssemblerNode()
	var port1 = NewPort()
	var port2 = NewPort()
	var port3 = NewPort()

	var pairPort1 = node.GetPairPort(port1)
	node.GetPairPort(port2)
	var pairPort3 = node.GetPairPort(port3)

	node.DeletePairPort(port2)

	var requirePorts = node.GetRequirePorts()

	assert.Equal(t, 2, len(requirePorts))
	assert.Equal(t, pairPort1, requirePorts[0])
	assert.Equal(t, pairPort3, requirePorts[1])
}
