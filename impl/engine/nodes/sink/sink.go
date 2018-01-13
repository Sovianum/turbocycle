package sink

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewSinkNode() SinkNode {
	var result = &sinkNode{}
	result.input = graph.NewAttachedPort(result)
	return result
}

func SinkAll(ports ...graph.Port) []SinkNode {
	var result = make([]SinkNode, len(ports))
	for i, port := range ports {
		result[i] = SinkPort(port)
	}
	return result
}

func SinkPort(port graph.Port) SinkNode {
	var node = NewSinkNode()
	graph.Link(node.Input(), port)
	return node
}

type SinkNode interface {
	graph.Node
	GetState() graph.PortState
	Input() graph.Port
}

type sinkNode struct {
	graph.BaseNode
	input graph.Port
}

func (node *sinkNode) GetState() graph.PortState {
	return node.input.GetState()
}

func (node *sinkNode) Input() graph.Port {
	return node.input
}

func (node *sinkNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Sink")
}

func (node *sinkNode) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{node.input}, nil
}

func (node *sinkNode) GetUpdatePorts() ([]graph.Port, error) {
	return make([]graph.Port, 0), nil
}

func (node *sinkNode) GetPorts() []graph.Port {
	return []graph.Port{node.input}
}

func (node *sinkNode) Process() error {
	return nil
}

func (node *sinkNode) setInput(input graph.Port) {
	node.input = input
}

func (node *sinkNode) getInput() graph.Port {
	return node.input
}
