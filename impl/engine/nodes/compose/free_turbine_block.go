package compose

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
)

type FreeTurbineBlockNode interface {
	graph.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
	FreeTurbine() constructive.FreeTurbineNode
	OutletPressureLoss() constructive.PressureLossNode
}

func NewFreeTurbineBlock(
	pAtm float64,
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(node constructive.TurbineNode) float64,
	pressureLossSigma float64,
) FreeTurbineBlockNode {
	var result = &freeTurbineBlockNode{
		// all values except pressure are sinked later
		atmNode: source.NewComplexGasSourceNode(nil, 0, pAtm, 1),
		turbine: constructive.NewFreeTurbineNode(
			etaT, lambdaOut, precision,
			leakMassRateFunc, coolMasRateRel, inflowMassRateRel,
		),
		pressureLoss: constructive.NewPressureLossNode(pressureLossSigma),
	}
	result.linkPorts()

	graph.AttachAllPorts(
		result,
		&result.powerOutput,
		&result.temperatureInput, &result.pressureInput, &result.gasInput, &result.massRateInput,
		&result.temperatureOutput, &result.pressureOutput, &result.gasOutput, &result.massRateOutput,
	)

	return result
}

type freeTurbineBlockNode struct {
	graph.BaseNode

	temperatureInput graph.Port
	pressureInput    graph.Port
	gasInput         graph.Port
	massRateInput    graph.Port

	temperatureOutput graph.Port
	pressureOutput    graph.Port
	gasOutput         graph.Port
	massRateOutput    graph.Port

	powerOutput graph.Port

	atmNode      source.ComplexGasSourceNode
	turbine      constructive.FreeTurbineNode
	pressureLoss constructive.PressureLossNode
	tSink        sink.SinkNode
	gSink        sink.SinkNode
	mSink        sink.SinkNode
}

func (node *freeTurbineBlockNode) GasInput() graph.Port {
	return node.gasInput
}

func (node *freeTurbineBlockNode) TemperatureInput() graph.Port {
	return node.temperatureInput
}

func (node *freeTurbineBlockNode) PressureInput() graph.Port {
	return node.pressureInput
}

func (node *freeTurbineBlockNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *freeTurbineBlockNode) GasOutput() graph.Port {
	return node.gasOutput
}

func (node *freeTurbineBlockNode) TemperatureOutput() graph.Port {
	return node.temperatureOutput
}

func (node *freeTurbineBlockNode) PressureOutput() graph.Port {
	return node.pressureOutput
}

func (node *freeTurbineBlockNode) MassRateOutput() graph.Port {
	return node.massRateOutput
}

func (node *freeTurbineBlockNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "FreeTurbineBlock")
}

func (node *freeTurbineBlockNode) GetPorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *freeTurbineBlockNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.temperatureInput, node.pressureInput, node.gasInput, node.massRateInput,
	}
}

func (node *freeTurbineBlockNode) GetUpdatePorts() []graph.Port {
	return []graph.Port{
		node.powerOutput,
		node.temperatureOutput, node.pressureOutput, node.gasOutput, node.massRateOutput,
	}
}

func (node *freeTurbineBlockNode) FreeTurbine() constructive.FreeTurbineNode {
	return node.turbine
}

func (node *freeTurbineBlockNode) OutletPressureLoss() constructive.PressureLossNode {
	return node.pressureLoss
}

func (node *freeTurbineBlockNode) Process() error {
	node.readInput()
	if err := node.atmNode.Process(); err != nil {
		return err
	}
	if err := node.pressureLoss.Process(); err != nil {
		return err
	}
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *freeTurbineBlockNode) PowerOutput() graph.Port {
	return node.powerOutput
}

func (node *freeTurbineBlockNode) linkPorts() {
	nodes.LinkComplexOutToOut(node.atmNode, node.pressureLoss)

	node.tSink = sink.SinkPort(node.pressureLoss.TemperatureInput())
	node.mSink = sink.SinkPort(node.pressureLoss.MassRateInput())
	node.gSink = sink.SinkPort(node.pressureLoss.GasInput())

	graph.Link(node.pressureLoss.PressureInput(), node.turbine.PressureOutput())
}

func (node *freeTurbineBlockNode) readInput() {
	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(), node.temperatureInput.GetState(),
			node.pressureInput.GetState(), node.massRateInput.GetState(),
		},
		[]graph.Port{
			node.turbine.GasInput(), node.turbine.TemperatureInput(),
			node.turbine.PressureInput(), node.turbine.MassRateInput(),
		},
	)
}

func (node *freeTurbineBlockNode) writeOutput() {
	graph.SetAll(
		[]graph.PortState{
			node.turbine.GasOutput().GetState(), node.turbine.TemperatureOutput().GetState(),
			node.turbine.PressureOutput().GetState(), node.turbine.MassRateOutput().GetState(),
		},
		[]graph.Port{
			node.gasOutput, node.temperatureOutput,
			node.pressureOutput, node.massRateOutput,
		},
	)
	node.powerOutput.SetState(node.turbine.PowerOutput().GetState())
}
