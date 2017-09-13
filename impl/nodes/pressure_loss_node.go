package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

type PressureLossNode interface {
	core.Node
	GasChannel
}

type pressureLossNode struct {
	ports core.PortsType
	sigma float64
}

func NewPressureLossNode(sigma float64) PressureLossNode {
	var result = &pressureLossNode{
		ports: make(core.PortsType),
		sigma: sigma,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasOutput] = core.NewPort()

	return result
}

func (node *pressureLossNode) ContextDefined() bool {
	return false	// TODO add real tests
}

func (node *pressureLossNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *pressureLossNode) Process() error {
	var inputState = node.gasInput().GetState()
	var outputState = node.gasOutput().GetState()

	if inputState == nil && outputState == nil {
		return errors.New("PressureLossNode failed to update any port (both states are nil)")
	}

	if inputState != nil && outputState == nil {
		var newOutputState = inputState.(states.GasPortState)
		var pStagOut = newOutputState.PStag * node.sigma
		newOutputState.PStag = pStagOut
		node.gasOutput().SetState(newOutputState)
		return nil
	}

	if inputState == nil && outputState != nil {
		var newInputState = outputState.(states.GasPortState)
		var pStagIn = newInputState.PStag / node.sigma
		newInputState.PStag = pStagIn
		node.GasInput().SetState(outputState)
		return nil
	}

	var newInputState, newOutputState = inputState.(states.GasPortState), outputState.(states.GasPortState)

	var pStagIn = newOutputState.PStag / node.sigma
	var pStagOut = newInputState.PStag * node.sigma

	newInputState.PStag = pStagIn
	newOutputState.PStag = pStagOut

	node.gasInput().SetState(newInputState)
	node.gasOutput().SetState(newOutputState)

	return nil
}

func (node *pressureLossNode) GetRequirePortTags() ([]string, error) {
	return []string{}, nil // TODO maybe need to make two modes; add context check
}

func (node *pressureLossNode) GetUpdatePortTags() ([]string, error) {
	return []string{gasInput, gasOutput}, nil	// TODO add context check
}

func (node *pressureLossNode) GetPortTags() []string {
	return []string{gasInput, gasOutput}
}

func (node *pressureLossNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case gasInput:
		return node.gasInput(), nil
	case gasOutput:
		return node.gasOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Failed to find port with tag \"%s\" in pressureLossNode", tag))
	}
}

func (node *pressureLossNode) GasOutput() core.Port {
	return node.gasOutput()
}

func (node *pressureLossNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *pressureLossNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *pressureLossNode) GasInput() core.Port {
	return node.gasInput()
}

func (node *pressureLossNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *pressureLossNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *pressureLossNode) gasInput() core.Port {
	return node.ports[gasInput]
}

func (node *pressureLossNode) gasOutput() core.Port {
	return node.ports[gasOutput]
}

func (node *pressureLossNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *pressureLossNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *pressureLossNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *pressureLossNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}
