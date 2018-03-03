package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type TemperatureSourceNode interface {
	graph.Node
	nodes.TemperatureSource
	SetTemperature(tStag float64)
	GetTemperature() float64
}

func NewTemperatureSourceNode(tStag float64) TemperatureSourceNode {
	var result = &temperatureSourceNode{
		tStag: tStag,
	}

	result.output = graph.NewAttachedPort(result)
	return result
}

type temperatureSourceNode struct {
	sourceNode
	tStag float64
}

func (node *temperatureSourceNode) GetTemperature() float64 {
	return node.tStag
}

func (node *temperatureSourceNode) SetTemperature(tStag float64) {
	node.tStag = tStag
}

func (node *temperatureSourceNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "TemperatureSource")
}

func (node *temperatureSourceNode) Process() error {
	node.output.SetState(states.NewTemperaturePortState(node.tStag))
	return nil
}

func (node *temperatureSourceNode) TStagOut() float64 {
	return node.tStag
}

func (node *temperatureSourceNode) TemperatureOutput() graph.Port {
	return node.output
}
