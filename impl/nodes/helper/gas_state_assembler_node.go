package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/nodes"
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

	result.ports[nodes.PressurePort] = core.NewPort()
	result.ports[nodes.PressurePort].SetInnerNode(result)
	result.ports[nodes.PressurePort].SetState(states.NewPressurePortState(1e5)) // TODO remove hardcode

	result.ports[nodes.TemperaturePort] = core.NewPort()
	result.ports[nodes.TemperaturePort].SetInnerNode(result)
	result.ports[nodes.TemperaturePort].SetState(states.NewTemperaturePortState(288)) // TODO remove hardcode

	result.ports[nodes.GasPort] = core.NewPort()
	result.ports[nodes.GasPort].SetInnerNode(result)
	result.ports[nodes.GasPort].SetState(states.NewGasPortState(gases.GetAir())) // TODO remove hardcode

	result.ports[nodes.ComplexGasPort] = core.NewPort()
	result.ports[nodes.ComplexGasPort].SetInnerNode(result)
	result.ports[nodes.ComplexGasPort].SetState(states.StandardAtmosphereState()) // TODO remove hardcode

	result.ports[nodes.MassRateRelPort] = core.NewPort()
	result.ports[nodes.MassRateRelPort].SetInnerNode(result)
	result.ports[nodes.MassRateRelPort].SetState(states.NewMassRateRelPortState(1)) // TODO remove hardcode

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
	var complexIsSource, complexErr = nodes.IsDataSource(node.ComplexGasPort())
	if complexErr != nil {
		return complexErr
	}

	if !complexIsSource {
		node.ComplexGasPort().SetState(states.NewComplexGasPortState(
			node.GasPort().GetState().(states.GasPortState).Gas,
			node.TemperaturePort().GetState().(states.TemperaturePortState).TStag,
			node.PressurePort().GetState().(states.PressurePortState).PStag,
			node.MassRateRelPort().GetState().(states.MassRateRelPortState).MassRateRel,
		))
	} else {
		var complexGasState = node.ComplexGasPort().GetState().(states.ComplexGasPortState)

		var gasIsSource, gasErr = nodes.IsDataSource(node.GasPort())
		if gasErr != nil {
			return gasErr
		}
		var pressureIsSource, pressureErr = nodes.IsDataSource(node.PressurePort())
		if pressureErr != nil {
			return pressureErr
		}
		var temperatureIsSource, temperatureErr = nodes.IsDataSource(node.TemperaturePort())
		if temperatureErr != nil {
			return temperatureErr
		}
		var massRateIsSource, massRateErr = nodes.IsDataSource(node.MassRateRelPort())
		if massRateErr != nil {
			return massRateErr
		}

		if !gasIsSource {
			node.GasPort().SetState(states.NewGasPortState(complexGasState.Gas))
		}
		if !pressureIsSource {
			node.PressurePort().SetState(states.NewPressurePortState(complexGasState.PStag))
		}
		if !temperatureIsSource {
			node.TemperaturePort().SetState(states.NewTemperaturePortState(complexGasState.TStag))
		}
		if !massRateIsSource {
			node.MassRateRelPort().SetState(states.NewMassRateRelPortState(complexGasState.MassRateRel))
		}
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
	return node.ports[nodes.ComplexGasPort]
}

func (node *gasStateAssemblerNode) PressurePort() core.Port {
	return node.ports[nodes.PressurePort]
}

func (node *gasStateAssemblerNode) TemperaturePort() core.Port {
	return node.ports[nodes.TemperaturePort]
}

func (node *gasStateAssemblerNode) MassRateRelPort() core.Port {
	return node.ports[nodes.MassRateRelPort]
}

func (node *gasStateAssemblerNode) GasPort() core.Port {
	return node.ports[nodes.GasPort]
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

	node.requirePortTags = requirePortTags
	return node.requirePortTags, nil
}

func (node *gasStateAssemblerNode) getUpdatePortTags() ([]string, error) {
	if node.updatePortTags != nil {
		return node.updatePortTags, nil
	}

	var updatePortTags, err = node.getPortTagsTemplate(func(isSource bool) bool {
		return !isSource
	})

	if err != nil {
		return nil, err
	}

	node.updatePortTags = updatePortTags
	return node.updatePortTags, nil
}

func (node *gasStateAssemblerNode) getPortTagsTemplate(condition func(isSource bool) bool) ([]string, error) {
	var portTags = make([]string, 0)

	for _, tag := range node.getPortTags() {
		var port, _ = node.getPortByTag(tag)
		var isSource, err = nodes.IsDataSource(port)
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
		var isSource, err = nodes.IsDataSource(port)
		if err != nil {
			return false, err
		}
		defined = defined && isSource
	}
	return defined, nil
}

func (node *gasStateAssemblerNode) getPortTags() []string {
	return []string{
		nodes.ComplexGasPort,
		nodes.GasPort,
		nodes.TemperaturePort,
		nodes.PressurePort,
		nodes.MassRateRelPort,
	}
}

func (node *gasStateAssemblerNode) getPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasPort:
		return node.ports[nodes.ComplexGasPort], nil
	case nodes.GasPort:
		return node.ports[nodes.GasPort], nil
	case nodes.TemperaturePort:
		return node.ports[nodes.TemperaturePort], nil
	case nodes.PressurePort:
		return node.ports[nodes.PressurePort], nil
	case nodes.MassRateRelPort:
		return node.ports[nodes.MassRateRelPort], nil
	default:
		return nil, errors.New(fmt.Sprintf("Port %s was not found in gasAssemblerNode", tag))
	}
}
