package constructive

import (
	"errors"
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

const (
	pressureNodeNotContextDefined = "pressure node not context defined"

	pressureLossInflow  = "pressureLossInflow"
	pressureLossOutflow = "pressureLossOutflow"
	pressureLossInitial = "pressureLossInitial"
)

type PressureLossNode interface {
	graph.Node
	nodes.ComplexGasChannel
	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut
	Sigma() float64
}

func NewPressureLossNode(sigma float64) PressureLossNode {
	var result = &pressureLossNode{
		sigma:          sigma,
		mode:           pressureLossInitial,
		contextCalled:  false,
		contextDefined: false,
	}

	result.complexGasInput = graph.NewAttachedPort(result)
	result.complexGasOutput = graph.NewAttachedPort(result)

	return result
}

type pressureLossNode struct {
	graph.BaseNode

	complexGasInput  graph.Port
	complexGasOutput graph.Port

	sigma          float64
	mode           string
	contextCalled  bool
	contextDefined bool
}

func (node *pressureLossNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PressureLossNode")
}

func (node *pressureLossNode) GetPorts() []graph.Port {
	return []graph.Port{node.complexGasInput, node.complexGasOutput}
}

func (node *pressureLossNode) GetRequirePorts() []graph.Port {
	var mode, err = node.getMode()
	if err != nil {
		panic(err)
	}

	if mode == pressureLossInflow {
		return []graph.Port{node.complexGasInput}
	}
	return []graph.Port{node.complexGasOutput}
}

func (node *pressureLossNode) GetUpdatePorts() []graph.Port {
	var mode, err = node.getMode()
	if err != nil {
		panic(err)
	}

	if mode == pressureLossOutflow {
		return []graph.Port{node.complexGasInput}
	}
	return []graph.Port{node.complexGasOutput}
}

func (node *pressureLossNode) ContextDefined() bool {
	if node.contextCalled {
		return node.contextDefined
	}
	node.contextCalled = true

	var node1 = node.complexGasInput.GetOuterNode()
	var node2 = node.complexGasOutput.GetOuterNode()

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

func (node *pressureLossNode) Process() error {
	var mode, contextErr = node.getMode()
	if contextErr != nil {
		return contextErr
	}

	switch mode {
	case pressureLossInflow:
		var inputState = node.complexGasInput.GetState()
		if inputState == nil {
			return errors.New("input state is nil")
		}

		var gasState = inputState.(states.ComplexGasPortState)
		gasState.PStag *= node.sigma
		node.complexGasOutput.SetState(gasState)
		return nil
	case pressureLossOutflow:
		var outputState = node.complexGasOutput.GetState()
		if outputState == nil {
			return errors.New("output state is nil")
		}

		var gasState = outputState.(states.ComplexGasPortState)
		gasState.PStag /= node.sigma
		node.complexGasInput.SetState(gasState)
		return nil
	default:
		return errors.New(pressureNodeNotContextDefined)
	}
}

func (node *pressureLossNode) ComplexGasOutput() graph.Port {
	return node.complexGasOutput
}

func (node *pressureLossNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *pressureLossNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *pressureLossNode) ComplexGasInput() graph.Port {
	return node.complexGasInput
}

func (node *pressureLossNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *pressureLossNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *pressureLossNode) Sigma() float64 {
	return node.sigma
}

func (node *pressureLossNode) tStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).TStag
}

func (node *pressureLossNode) pStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).PStag
}

func (node *pressureLossNode) tStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *pressureLossNode) pStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).PStag
}

func (node *pressureLossNode) getMode() (string, error) {
	if node.mode != pressureLossInitial {
		return node.mode, nil
	}

	var inputIsSource, inputErr = nodes.IsDataSource(node.complexGasInput)
	if inputErr != nil {
		return "", inputErr
	}
	var outputIsSource, outputErr = nodes.IsDataSource(node.complexGasOutput)
	if outputErr != nil {
		return "", outputErr
	}
	if inputIsSource && !outputIsSource {
		return pressureLossInflow, nil
	}
	if !inputIsSource && outputIsSource {
		return pressureLossOutflow, nil
	}
	return "", fmt.Errorf("inconsistent pressure loss node state")
}
