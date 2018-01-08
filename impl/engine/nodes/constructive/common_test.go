package constructive

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiAdder(t *testing.T) {
	var port1, port2, port3, port4 = graph.NewPort(), graph.NewPort(), graph.NewPort(), graph.NewPort()

	port1.SetState(graph.NewNumberPortState(2))
	port2.SetState(graph.NewNumberPortState(3))
	port3.SetState(graph.NewNumberPortState(4))
	port4.SetState(graph.NewNumberPortState(5))

	var node = NewMultiAdder()
	node.AddPortGroup(port1, port2)
	node.AddPortGroup(port3, port4)

	var err = node.Process()
	assert.Nil(t, err)

	assert.InDelta(t, 26, node.OutputPort().GetState().Value().(float64), 1e-7)
}

func TestNewEquality(t *testing.T) {
	var port1, port2 = graph.NewPort(), graph.NewPort()

	port1.SetState(graph.NewNumberPortState(2))
	port2.SetState(graph.NewNumberPortState(3))

	var node = NewEquality(port1, port2)

	var err = node.Process()
	assert.Nil(t, err)

	assert.InDelta(t, -1, node.OutputPort().GetState().Value().(float64), 1e-7)
}
