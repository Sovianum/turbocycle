package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type TemperatureSinkNode interface {
	graph.Node
	nodes.TemperatureSink
}

func NewTemperatureSinkNode() TemperatureSinkNode {
	var result = &temperatureSinkNode{}
	result.setInput(graph.NewAttachedPort(result))
	return result
}

type temperatureSinkNode struct {
	sinkNode
}

func (node *temperatureSinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TemperatureSink")
}

func (node *temperatureSinkNode) TStagIn() float64 {
	return node.getInput().GetState().(states.TemperaturePortState).TStag
}

func (node *temperatureSinkNode) TemperatureInput() graph.Port {
	return node.getInput()
}
