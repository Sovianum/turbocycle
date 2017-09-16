package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/states"
)

const (
	pressureNodeNotContextDefined = "Pressure node not context defined"

	pressureLossInflow  = "pressureLossInflow"
	pressureLossOutflow = "pressureLossOutflow"
	pressureLossBiFlow  = "pressureLossBiFlow"
	pressureLossInitial = "pressureLossInitial"
)

type PressureLossNode interface {
	core.Node
	GasChannel
}

type pressureLossNode struct {
	ports          core.PortsType
	sigma          float64
	mode           string
	contextCalled  bool
	contextDefined bool
}

func NewPressureLossNode(sigma float64) PressureLossNode {
	var result = &pressureLossNode{
		ports:          make(core.PortsType),
		sigma:          sigma,
		mode:           pressureLossInitial,
		contextCalled:  false,
		contextDefined: false,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetInnerNode(result)
	result.ports[gasInput].SetState(states.StandardAtmosphereState())

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)
	result.ports[gasOutput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *pressureLossNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState  core.PortState `json:"gas_input_state"`
		GasOutputState core.PortState `json:"gas_output_state"`
		Sigma          float64        `json:"sigma"`
	}{
		GasInputState:  node.gasInput().GetState(),
		GasOutputState: node.gasOutput().GetState(),
		Sigma:          node.sigma,
	})
}

func (node *pressureLossNode) ContextDefined() bool {
	if node.contextCalled {
		return node.contextDefined
	}
	node.contextCalled = true

	var node1 = node.gasInput().GetOuterNode()
	var node2 = node.gasOutput().GetOuterNode()

	if node1 != nil {
		var node1Defined = node1.ContextDefined()
		if node1Defined {
			node.contextDefined = true
			return node.contextDefined
		}
	}

	if node2 != nil {
		var node2Defined = node2.ContextDefined()
		if node2Defined {
			node.contextDefined = true
			return node.contextDefined
		}
	}

	node.contextDefined = false
	return node.contextDefined
}

func (node *pressureLossNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *pressureLossNode) Process() error {
	var mode, contextErr = node.getMode()
	if contextErr != nil {
		return contextErr
	}

	switch mode {
	case pressureLossInflow:
		var inputState = node.gasInput().GetState()
		if inputState == nil {
			return errors.New("Input state is nil")
		}

		var gasState = inputState.(states.GasPortState)
		gasState.PStag *= node.sigma
		node.gasOutput().SetState(gasState)
		return nil
	case pressureLossOutflow:
		var outputState = node.gasOutput().GetState()
		if outputState == nil {
			return errors.New("Output state is nil")
		}

		var gasState = outputState.(states.GasPortState)
		gasState.PStag /= node.sigma
		node.gasInput().SetState(gasState)
		return nil
	case pressureLossBiFlow:
		var inputState = node.gasInput().GetState()
		var outputState = node.gasOutput().GetState()
		if inputState == nil && outputState == nil {
			return errors.New("Both input and output states are nil")
		}

		if inputState != nil && outputState == nil {
			var gasState = inputState.(states.GasPortState)
			gasState.PStag *= node.sigma
			node.gasOutput().SetState(gasState)
			return nil
		}

		if inputState == nil && outputState != nil {
			var gasState = outputState.(states.GasPortState)
			gasState.PStag /= node.sigma
			node.gasInput().SetState(gasState)
			return nil
		}

		if inputState != nil && outputState != nil {
			var inputGasState = inputState.(states.GasPortState)
			var outputGasState = outputState.(states.GasPortState)

			inputGasState.PStag, outputGasState.PStag = outputGasState.PStag/node.sigma, inputGasState.PStag*node.sigma

			node.gasInput().SetState(inputGasState)
			node.gasOutput().SetState(outputGasState)
		}
		return nil
	default:
		return errors.New(pressureNodeNotContextDefined)
	}
}

func (node *pressureLossNode) GetRequirePortTags() ([]string, error) {
	var mode, err = node.getMode()
	if err != nil {
		return nil, err
	}
	switch mode {
	case pressureLossInflow:
		return []string{gasInput}, nil
	case pressureLossOutflow:
		return []string{gasOutput}, nil
	case pressureLossBiFlow:
		return []string{}, nil
	default:
		return nil, errors.New(pressureNodeNotContextDefined)
	}
}

func (node *pressureLossNode) GetUpdatePortTags() ([]string, error) {
	var mode, err = node.getMode()
	if err != nil {
		return nil, err
	}
	switch mode {
	case pressureLossInflow:
		return []string{gasOutput}, nil
	case pressureLossOutflow:
		return []string{gasInput}, nil
	case pressureLossBiFlow:
		return []string{gasInput, gasOutput}, nil
	default:
		return nil, errors.New(pressureNodeNotContextDefined)
	}
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

func (node *pressureLossNode) getMode() (string, error) {
	if node.mode != pressureLossInitial {
		return node.mode, nil
	}

	var inputIsSource, inputErr = isDataSource(node.gasInput())
	if inputErr != nil {
		return "", inputErr
	}
	var outputIsSource, outputErr = isDataSource(node.gasOutput())
	if outputErr != nil {
		return "", outputErr
	}
	if inputIsSource && outputIsSource {
		return pressureLossBiFlow, nil
	}
	if inputIsSource {
		return pressureLossInflow, nil
	}
	if outputIsSource {
		return pressureLossOutflow, nil
	}
	return pressureLossInitial, nil
}

func isDataSource(port core.Port) (bool, error) {
	var linkPort = port.GetLinkPort()
	if linkPort == nil {
		return false, nil
	}

	var outerNode = port.GetOuterNode()
	if outerNode == nil {
		return false, nil
	}

	if !outerNode.ContextDefined() {
		return false, nil
	}

	var updatePortTags, err = outerNode.GetUpdatePortTags()
	if err != nil {
		return false, err
	}

	for _, updatePortTag := range updatePortTags {
		var tagPort, tagErr = outerNode.GetPortByTag(updatePortTag)
		if tagErr != nil {
			return false, tagErr
		}

		if tagPort == linkPort {
			return true, nil
		}
	}

	return false, nil
}
