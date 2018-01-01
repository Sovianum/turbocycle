package source

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type TemperatureSourceNode interface {
	core.Node
	nodes.TemperatureSource
}

func NewTemperatureSourceNode(tStag float64) TemperatureSourceNode {
	var result = &temperatureSourceNode{
		tStag: tStag,
	}

	result.output = core.NewAttachedPort(result)
	return result
}

type temperatureSourceNode struct {
	sourceNode
	tStag float64
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

func (node *temperatureSourceNode) TemperatureOutput() core.Port {
	return node.output
}
