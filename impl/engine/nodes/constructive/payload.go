package constructive

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type Payload interface {
	graph.Node
	nodes.RPMSource
	nodes.PowerSource

	NormRPM() float64
	SetNormRPM(normRPM float64)

	RPM() float64
	Power() float64
}

func NewPayload(rpm0, power0 float64, powerCharacteristic func(normRpm float64) float64) Payload {
	var result = &payload{
		powerCharacteristic: powerCharacteristic,

		rpm0:    rpm0,
		power0:  power0,
		normRpm: 1,
	}
	graph.AttachAllWithTags(
		result,
		[]*graph.Port{
			&result.rpmOutput, &result.powerOutput,
		},
		[]string{
			nodes.RPMOutputTag, nodes.PowerOutputTag,
		},
	)
	return result
}

type payload struct {
	graph.BaseNode

	rpmOutput   graph.Port
	powerOutput graph.Port

	rpm0   float64
	power0 float64

	powerCharacteristic func(normRpm float64) float64

	normRpm float64
}

func (node *payload) NormRPM() float64 {
	return node.normRpm
}

func (node *payload) SetNormRPM(normRPM float64) {
	node.normRpm = normRPM
}

func (node *payload) RPM() float64 {
	return node.rpm0 * node.normRpm
}

func (node *payload) Power() float64 {
	return node.powerCharacteristic(node.normRpm) * node.power0
}

func (node *payload) GetName() string {
	return common.EitherString(node.GetInstanceName(), "Payload")
}

func (node *payload) Process() error {
	var rpm = node.normRpm * node.rpm0
	var power = node.powerCharacteristic(node.normRpm) * node.power0

	node.rpmOutput.SetState(states.NewRPMPortState(rpm))
	node.powerOutput.SetState(states.NewPowerPortState(power))

	return nil
}

func (node *payload) GetRequirePorts() ([]graph.Port, error) {
	return make([]graph.Port, 0), nil
}

func (node *payload) GetUpdatePorts() ([]graph.Port, error) {
	return []graph.Port{node.rpmOutput, node.powerOutput}, nil
}

func (node *payload) GetPorts() []graph.Port {
	return []graph.Port{node.rpmOutput, node.powerOutput}
}

func (node *payload) RPMOutput() graph.Port {
	return node.rpmOutput
}

func (node *payload) PowerOutput() graph.Port {
	return node.powerOutput
}
