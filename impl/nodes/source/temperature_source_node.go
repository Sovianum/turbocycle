package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/states"
)

type TemperatureSourceNode interface {
	core.Node
	nodes.TemperatureSource
}

type temperatureSourceNode struct {
	ports core.PortsType
	tStag float64
}

func NewTemperatureSourceNode(tStag float64) TemperatureSourceNode {
	var result = &temperatureSourceNode{
		ports: make(core.PortsType),
		tStag: tStag,
	}

	result.ports[nodes.TemperatureOutput] = core.NewPort()
	result.ports[nodes.TemperatureOutput].SetInnerNode(result)
	result.ports[nodes.TemperatureOutput].SetState(states.NewTemperaturePortState(tStag))

	return result
}

func (node *temperatureSourceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TemperatureOutputState core.PortState `json:"temperature_output_state"`
	}{
		TemperatureOutputState: node.TemperatureOutput().GetState(),
	})
}

func (node *temperatureSourceNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *temperatureSourceNode) Process() error {
	node.ports[nodes.TemperatureOutput].SetState(states.NewTemperaturePortState(node.tStag))
	return nil
}

func (node *temperatureSourceNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil
}

func (node *temperatureSourceNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.TemperatureOutput}, nil
}

func (node *temperatureSourceNode) GetPortTags() []string {
	return []string{nodes.TemperatureOutput}
}

func (node *temperatureSourceNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.TemperatureOutput:
		return node.ports[nodes.TemperatureOutput], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s of temperatureSourceNode can not be found", tag))
	}
}

func (node *temperatureSourceNode) ContextDefined() bool {
	return true
}

func (node *temperatureSourceNode) TStagOut() float64 {
	return node.ports[nodes.TemperatureOutput].GetState().(states.TemperaturePortState).TStag
}

func (node *temperatureSourceNode) TemperatureOutput() core.Port {
	return node.ports[nodes.TemperatureOutput]
}
