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

	nodes.GasChannel
	nodes.TemperatureChannel
	nodes.PressureChannel
	nodes.MassRateChannel

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

	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
			&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
		},
		[]string{
			nodes.TemperatureInputTag, nodes.PressureInputTag, nodes.GasInputTag, nodes.MassRateInputTag,
			nodes.TemperatureOutputTag, nodes.PressureOutputTag, nodes.GasOutputTag, nodes.MassRateOutputTag,
		},
	)

	return result
}

type pressureLossNode struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	pressureOutput    graph.Port
	temperatureOutput graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	sigma          float64
	mode           string
	contextCalled  bool
	contextDefined bool
}

func (node *pressureLossNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *pressureLossNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *pressureLossNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *pressureLossNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *pressureLossNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *pressureLossNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *pressureLossNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *pressureLossNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *pressureLossNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "PressureLossNode")
}

func (node *pressureLossNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *pressureLossNode) GetRequirePorts() ([]graph.Port, error) {
	var mode, err = node.getMode()
	if err != nil {
		return nil, err
	}

	if mode == pressureLossInflow {
		return []graph.Port{
			node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		}, nil
	}
	return []graph.Port{
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}, nil
}

func (node *pressureLossNode) GetUpdatePorts() ([]graph.Port, error) {
	var mode, err = node.getMode()
	if err != nil {
		return nil, err
	}

	if mode == pressureLossOutflow {
		return []graph.Port{
			node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		}, nil
	}
	return []graph.Port{
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}, nil
}

func (node *pressureLossNode) ContextDefined(key int) bool {
	if node.contextCalled {
		return node.contextDefined
	}
	node.contextCalled = true

	var node1 = node.gasInput.GetOuterNode()
	var node2 = node.gasOutput.GetOuterNode()

	if node1 != nil {
		var node1Defined = node1.ContextDefined(0)
		if node1Defined {
			node.contextDefined = true
			return node.contextDefined
		}
	}

	if node2 != nil {
		var node2Defined = node2.ContextDefined(0)
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
		graph.SetAll(
			[]graph.PortState{
				node.gasInput.GetState(),
				node.temperatureInput.GetState(),
				states.NewPressurePortState(node.pStagIn() * node.sigma),
				node.massRateInput.GetState(),
			},
			[]graph.Port{node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput},
		)

		return nil
	case pressureLossOutflow:
		graph.SetAll(
			[]graph.PortState{
				node.gasOutput.GetState(),
				node.temperatureOutput.GetState(),
				states.NewPressurePortState(node.pStagOut() / node.sigma),
				node.massRateOutput.GetState(),
			},
			[]graph.Port{node.gasInput, node.temperatureInput, node.pressureInput, node.massRateInput},
		)

		return nil
	default:
		return errors.New(pressureNodeNotContextDefined)
	}
}

func (node *pressureLossNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *pressureLossNode) PStagOut() float64 {
	return node.pStagOut()
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
	return node.temperatureOutput.GetState().(states.TemperaturePortState).TStag
}

func (node *pressureLossNode) pStagOut() float64 {
	return node.pressureOutput.GetState().(states.PressurePortState).PStag
}

func (node *pressureLossNode) tStagIn() float64 {
	return node.temperatureInput.GetState().(states.TemperaturePortState).TStag
}

func (node *pressureLossNode) pStagIn() float64 {
	return node.pressureInput.GetState().(states.PressurePortState).PStag
}

func (node *pressureLossNode) getMode() (string, error) {
	if node.mode != pressureLossInitial {
		return node.mode, nil
	}

	var inputIsSource, inputErr = nodes.IsDataSource(node.pressureInput)
	if inputErr != nil {
		return "", inputErr
	}
	var outputIsSource, outputErr = nodes.IsDataSource(node.pressureOutput)
	if outputErr != nil {
		return "", outputErr
	}
	if inputIsSource && !outputIsSource {
		return pressureLossInflow, nil
	}
	if !inputIsSource && outputIsSource {
		return pressureLossOutflow, nil
	}
	if !inputIsSource && !outputIsSource {
		return "", fmt.Errorf("neither endpoint is pressure source")
	}
	return "", fmt.Errorf("both endpoints are pressure sources")
}
