package common

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/stage/states"
)

func LinkStages(source, dest StageChannel) {
	graph.LinkAll(
		[]graph.Port{
			source.GasOutput(), source.PressureOutput(),
			source.TemperatureOutput(), source.MassRateOutput(),
			source.VelocityOutput(),
		},
		[]graph.Port{
			dest.GasInput(), dest.PressureInput(),
			dest.TemperatureInput(), dest.MassRateInput(),
			dest.VelocityInput(),
		},
	)
}

func InitFromPreviousStage(source, dest StageChannel) {
	graph.CopyAll(
		[]graph.Port{
			source.GasOutput(), source.PressureOutput(),
			source.TemperatureOutput(), source.MassRateOutput(),
			source.VelocityOutput(),
		},
		[]graph.Port{
			dest.GasInput(), dest.PressureInput(),
			dest.TemperatureInput(), dest.MassRateInput(),
			dest.VelocityInput(),
		},
	)
}

type StageChannel interface {
	graph.Node
	nodes.GasChannel
	nodes.PressureChannel
	nodes.TemperatureChannel
	VelocityChannel
	nodes.MassRateChannel
}

func NewBaseStage(node graph.Node) *BaseStage {
	result := new(BaseStage)
	result.attachPorts(node)
	return result
}

type BaseStage struct {
	graph.BaseNode

	gasInput         graph.Port
	temperatureInput graph.Port
	pressureInput    graph.Port
	massRateInput    graph.Port
	velocityInput    graph.Port

	gasOutput         graph.Port
	temperatureOutput graph.Port
	pressureOutput    graph.Port
	massRateOutput    graph.Port
	velocityOutput    graph.Port
}

func (base *BaseStage) attachPorts(node graph.Node) {
	graph.AttachAllWithTags(
		node,
		[]*graph.Port{
			&base.gasInput, &base.gasOutput,
			&base.pressureInput, &base.pressureOutput,
			&base.temperatureInput, &base.temperatureOutput,
			&base.velocityInput, &base.velocityOutput,
			&base.massRateInput, &base.massRateOutput,
		},
		[]string{
			nodes.GasInputTag, nodes.GasOutputTag,
			nodes.PressureInputTag, nodes.PressureOutputTag,
			nodes.TemperatureInputTag, nodes.TemperatureOutputTag,
			states.VelocityInletTag, states.VelocityOutletTag,
			nodes.MassRateInputTag, nodes.MassRateOutputTag,
		},
	)
}

func (base *BaseStage) GetRequirePorts() ([]graph.Port, error) {
	return []graph.Port{base.gasInput, base.temperatureInput, base.pressureInput, base.massRateInput, base.velocityInput}, nil
}

func (base *BaseStage) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{base.gasOutput, base.temperatureOutput, base.pressureOutput, base.massRateOutput, base.velocityOutput}, nil
}

func (base *BaseStage) GetPorts() []graph.Port {
	return []graph.Port{
		base.gasInput, base.temperatureInput, base.pressureInput, base.massRateInput, base.velocityInput,
		base.gasOutput, base.temperatureOutput, base.pressureOutput, base.massRateOutput, base.velocityOutput,
	}
}

func (base *BaseStage) GasOutput() graph.Port {
	return base.gasOutput
}

func (base *BaseStage) GasInput() graph.Port {
	return base.gasInput
}

func (base *BaseStage) PressureOutput() graph.Port {
	return base.pressureOutput
}

func (base *BaseStage) PressureInput() graph.Port {
	return base.pressureInput
}

func (base *BaseStage) TemperatureOutput() graph.Port {
	return base.temperatureOutput
}

func (base *BaseStage) TemperatureInput() graph.Port {
	return base.temperatureInput
}

func (base *BaseStage) VelocityInput() graph.Port {
	return base.velocityInput
}

func (base *BaseStage) VelocityOutput() graph.Port {
	return base.velocityOutput
}

func (base *BaseStage) MassRateInput() graph.Port {
	return base.massRateInput
}

func (base *BaseStage) MassRateOutput() graph.Port {
	return base.massRateOutput
}
