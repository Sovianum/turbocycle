package helper

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type GasStateAssemblerNode interface {
	graph.Node
	ComplexGasPort() graph.Port
	PressurePort() graph.Port
	TemperaturePort() graph.Port
	MassRateRelPort() graph.Port
	GasPort() graph.Port
}

func NewGasStateDisassemblerNode() GasStateAssemblerNode {
	return newAssemblerNode(false)
}

func NewGasStateAssemblerNode() GasStateAssemblerNode {
	return newAssemblerNode(true)
}

func newAssemblerNode(isAssembler bool) GasStateAssemblerNode {
	var result = &gasStateAssemblerNode{
		contextCalledFlag:  false,
		contextDefinedFlag: false,
		requirePortTags:    nil,
		updatePortTags:     nil,
		isAssembler:        isAssembler,
	}

	result.complexGasPort = graph.NewAttachedPort(result)
	result.pressurePort = graph.NewAttachedPort(result)
	result.temperaturePort = graph.NewAttachedPort(result)
	result.massRateRelPort = graph.NewAttachedPort(result)
	result.gasPort = graph.NewAttachedPort(result)

	return result
}

type gasStateAssemblerNode struct {
	graph.BaseNode

	complexGasPort  graph.Port
	pressurePort    graph.Port
	temperaturePort graph.Port
	massRateRelPort graph.Port
	gasPort         graph.Port

	contextCalledFlag  bool
	contextDefinedFlag bool
	requirePortTags    []string
	updatePortTags     []string
	isAssembler        bool
}

func (node *gasStateAssemblerNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "GasStateAssembler")
}

func (node *gasStateAssemblerNode) GetRequirePorts() []graph.Port {
	if node.isAssembler {
		return []graph.Port{
			node.gasPort,
			node.temperaturePort,
			node.pressurePort,
			node.massRateRelPort,
		}
	} else {
		return []graph.Port{
			node.complexGasPort,
		}
	}
}

func (node *gasStateAssemblerNode) GetUpdatePorts() []graph.Port {
	if !node.isAssembler {
		return []graph.Port{
			node.gasPort,
			node.temperaturePort,
			node.pressurePort,
			node.massRateRelPort,
		}
	} else {
		return []graph.Port{
			node.complexGasPort,
		}
	}
}

func (node *gasStateAssemblerNode) GetPorts() []graph.Port {
	return node.getPorts()
}

func (node *gasStateAssemblerNode) Process() error {
	if node.isAssembler {
		node.complexGasPort.SetState(states.NewComplexGasPortState(
			node.gasPort.GetState().(states.GasPortState).Gas,
			node.temperaturePort.GetState().(states.TemperaturePortState).TStag,
			node.pressurePort.GetState().(states.PressurePortState).PStag,
			node.massRateRelPort.GetState().(states.MassRateRelPortState).MassRateRel,
		))
	} else {
		var complexGasState = node.complexGasPort.GetState().(states.ComplexGasPortState)

		var gasIsSource, gasErr = nodes.IsDataSource(node.gasPort)
		if gasErr != nil {
			return gasErr
		}
		var pressureIsSource, pressureErr = nodes.IsDataSource(node.pressurePort)
		if pressureErr != nil {
			return pressureErr
		}
		var temperatureIsSource, temperatureErr = nodes.IsDataSource(node.temperaturePort)
		if temperatureErr != nil {
			return temperatureErr
		}
		var massRateIsSource, massRateErr = nodes.IsDataSource(node.massRateRelPort)
		if massRateErr != nil {
			return massRateErr
		}

		if !gasIsSource {
			node.gasPort.SetState(states.NewGasPortState(complexGasState.Gas))
		}
		if !pressureIsSource {
			node.pressurePort.SetState(states.NewPressurePortState(complexGasState.PStag))
		}
		if !temperatureIsSource {
			node.temperaturePort.SetState(states.NewTemperaturePortState(complexGasState.TStag))
		}
		if !massRateIsSource {
			node.massRateRelPort.SetState(states.NewMassRateRelPortState(complexGasState.MassRateRel))
		}
	}

	return nil
}

func (node *gasStateAssemblerNode) ContextDefined() bool {
	if node.contextCalledFlag {
		return node.contextDefinedFlag
	}
	node.contextCalledFlag = true

	var defined = true
	for _, port := range node.getPorts() {
		var outerNode = port.GetOuterNode()
		defined = defined && outerNode.ContextDefined()
	}

	node.contextDefinedFlag = defined
	node.contextCalledFlag = false
	return node.contextDefinedFlag
}

func (node *gasStateAssemblerNode) ComplexGasPort() graph.Port {
	return node.complexGasPort
}

func (node *gasStateAssemblerNode) PressurePort() graph.Port {
	return node.pressurePort
}

func (node *gasStateAssemblerNode) TemperaturePort() graph.Port {
	return node.temperaturePort
}

func (node *gasStateAssemblerNode) MassRateRelPort() graph.Port {
	return node.massRateRelPort
}

func (node *gasStateAssemblerNode) GasPort() graph.Port {
	return node.gasPort
}

func (node *gasStateAssemblerNode) getPorts() []graph.Port {
	return []graph.Port{
		node.complexGasPort,
		node.pressurePort,
		node.temperaturePort,
		node.massRateRelPort,
		node.gasPort,
	}
}
