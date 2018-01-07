package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorAssemblerNode_GetPairPort(t *testing.T) {
	var node = NewVectorAssemblerNode()
	var port1 = NewPort()
	var port2 = NewPort()

	node.AddInputPorts(port1)
	node.AddInputPorts(port2)

	var requirePorts = node.GetRequirePorts()

	assert.Equal(t, 2, len(requirePorts))
}

func TestVectorAssemblerNode_DeletePairPort(t *testing.T) {
	var node = NewVectorAssemblerNode()
	var port1 = NewPort()
	var port2 = NewPort()
	var port3 = NewPort()

	node.AddInputPorts(port1)
	node.AddInputPorts(port2)
	node.AddInputPorts(port3)

	node.DeleteInputPorts(port2)

	var requirePorts = node.GetRequirePorts()

	assert.Equal(t, 2, len(requirePorts))
}
