package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type GasStateAssemblerNode interface {
	core.Node
	ComplexGasPort() core.Port
	PressurePort() core.Port
	TemperaturePort() core.Port
	MassRateRelPort() core.Port
	GasPort() core.Port
}

type gasStateAssemblerNode struct {
	ports              core.PortsType
	contextCalledFlag  bool
	contextDefinedFlag bool
	requirePortTags    []string
	updatePortTags     []string
}

func NewGasStateAssemblerNode() GasStateAssemblerNode {
	var result = &gasStateAssemblerNode{
		ports:              make(core.PortsType),
		contextCalledFlag:  false,
		contextDefinedFlag: false,
		requirePortTags:    nil,
		updatePortTags:     nil,
	}

	result.ports[pressurePort] = core.NewPort()
	result.ports[pressurePort].SetInnerNode(result)

	result.ports[temperaturePort] = core.NewPort()
	result.ports[temperaturePort].SetInnerNode(result)

	result.ports[gasPort] = core.NewPort()
	result.ports[gasPort].SetInnerNode(result)

	result.ports[complexGasPort] = core.NewPort()
	result.ports[complexGasPort].SetInnerNode(result)

	result.ports[massRateRelPort] = core.NewPort()
	result.ports[massRateRelPort].SetInnerNode(result)

	return result
}

func (node *gasStateAssemblerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ComplexGasPortState  core.PortState
		PressurePortState    core.PortState
		TemperaturePortState core.PortState
		GasPortState         core.PortState
	}{
		ComplexGasPortState:  node.ComplexGasPort().GetState(),
		PressurePortState:    node.PressurePort().GetState(),
		TemperaturePortState: node.TemperaturePort().GetState(),
		GasPortState:         node.GasPort().GetState(),
	})
}

func (node *gasStateAssemblerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasStateAssemblerNode) Process() error {
	var complexIsSource, complexErr = isDataSource(node.ComplexGasPort())
	if complexErr != nil {
		return complexErr
	}

	if !complexIsSource {
		var complexGasState = node.ComplexGasPort().GetState().(states.ComplexGasPortState)

		var gasIsSource, gasErr = isDataSource(node.GasPort())
		if gasErr != nil {
			return gasErr
		}
		var pressureIsSource, pressureErr = isDataSource(node.PressurePort())
		if pressureErr != nil {
			return pressureErr
		}
		var temperatureIsSource, temperatureErr = isDataSource(node.TemperaturePort())
		if temperatureErr != nil {
			return temperatureErr
		}
		var massRateIsSource, massRateErr = isDataSource(node.MassRateRelPort())
		if massRateErr != nil {
			return massRateErr
		}

		if gasIsSource {
			node.GasPort().SetState(states.NewGasPortState(complexGasState.Gas))
		}
		if pressureIsSource {
			node.PressurePort().SetState(states.NewPressurePortState(complexGasState.PStag))
		}
		if temperatureIsSource {
			node.TemperaturePort().SetState(states.NewTemperaturePortState(complexGasState.TStag))
		}
		if massRateIsSource {
			node.MassRateRelPort().SetState(states.NewMassRateRelPortState(complexGasState.MassRateRel))
		}
	} else {
		node.ComplexGasPort().SetState(states.NewComplexGasPortState(
			node.GasPort().GetState().(states.GasPortState).Gas,
			node.TemperaturePort().GetState().(states.TemperaturePortState).TStag,
			node.PressurePort().GetState().(states.PressurePortState).PStag,
			node.MassRateRelPort().GetState().(states.MassRateRelPortState).MassRateRel,
		))
	}

	return nil
}

func (node *gasStateAssemblerNode) GetRequirePortTags() ([]string, error) {
	return node.getRequirePortTags()
}

func (node *gasStateAssemblerNode) GetUpdatePortTags() ([]string, error) {
	return node.getUpdatePortTags()
}

func (node *gasStateAssemblerNode) GetPortTags() []string {
	return node.getPortTags()
}

func (node *gasStateAssemblerNode) GetPortByTag(tag string) (core.Port, error) {
	return node.getPortByTag(tag)
}

func (node *gasStateAssemblerNode) ContextDefined() bool {
	var defined, err = node.contextDefined()
	if err != nil {
		return false
	}
	return defined
}

func (node *gasStateAssemblerNode) ComplexGasPort() core.Port {
	return node.ports[complexGasPort]
}

func (node *gasStateAssemblerNode) PressurePort() core.Port {
	return node.ports[pressurePort]
}

func (node *gasStateAssemblerNode) TemperaturePort() core.Port {
	return node.ports[temperaturePort]
}

func (node *gasStateAssemblerNode) MassRateRelPort() core.Port {
	return node.ports[massRateRelPort]
}

func (node *gasStateAssemblerNode) GasPort() core.Port {
	return node.ports[gasPort]
}

func (node *gasStateAssemblerNode) getRequirePortTags() ([]string, error) {
	if node.requirePortTags != nil {
		return node.requirePortTags, nil
	}

	var requirePortTags, err = node.getPortTagsTemplate(func(isSource bool) bool {
		return isSource
	})

	if err != nil {
		return nil, err
	}

	return requirePortTags, nil
}

func (node *gasStateAssemblerNode) getUpdatePortTags() ([]string, error) {
	if node.requirePortTags != nil {
		return node.requirePortTags, nil
	}

	var requirePortTags, err = node.getPortTagsTemplate(func(isSource bool) bool {
		return !isSource
	})

	if err != nil {
		return nil, err
	}

	return requirePortTags, nil
}

func (node *gasStateAssemblerNode) getPortTagsTemplate(condition func(isSource bool)bool) ([]string, error) {
	var portTags = make([]string, 0)

	for _, tag := range node.getPortTags() {
		var port, _ = node.getPortByTag(tag)
		var isSource, err = isDataSource(port)
		if err != nil {
			return nil, err
		}
		if condition(isSource) {
			portTags = append(portTags, tag)
		}
	}

	return portTags, nil
}

func (node *gasStateAssemblerNode) contextDefined() (bool, error) {
	if node.contextCalledFlag {
		return node.contextDefinedFlag, nil
	}
	node.contextCalledFlag = true

	var defined = true
	for _, tag := range node.getPortTags() {
		var port, _ = node.getPortByTag(tag)
		var isSource, err = isDataSource(port)
		if err != nil {
			return false, err
		}
		defined = defined && isSource
	}
	return defined, nil
}

func (node *gasStateAssemblerNode) getPortTags() []string {
	return []string{complexGasPort, gasPort, temperaturePort, pressurePort}
}

func (node *gasStateAssemblerNode) getPortByTag(tag string) (core.Port, error) {
	switch tag {
	case complexGasPort:
		return node.ports[complexGasPort], nil
	case gasPort:
		return node.ports[gasPort], nil
	case temperaturePort:
		return node.ports[temperaturePort], nil
	case pressurePort:
		return node.ports[pressurePort], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s was not found in gasAssemblerNode", tag))
	}
}
