package compose

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/sink"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func NewTurboCascadeNode(
	compressorEtaAd, piStag float64,
	etaT, lambdaOut float64, turbineMassRateRelFunc func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) TurboCascadeNode {
	var result = &turboCascadeNode{
		ports:        make(core.PortsType),
		compressor:   constructive.NewCompressorNode(compressorEtaAd, piStag, precision),
		turbine:      constructive.NewBlockedTurbineNode(etaT, lambdaOut, precision, turbineMassRateRelFunc),
		transmission: constructive.NewTransmissionNode(etaM),
		powerSink:    sink.NewPowerSinkNode(),
	}

	result.linkPorts()

	result.ports[nodes.CompressorComplexGasInput] = core.NewPort()
	result.ports[nodes.CompressorComplexGasInput].SetInnerNode(result)
	result.ports[nodes.CompressorComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.CompressorComplexGasOutput] = core.NewPort()
	result.ports[nodes.CompressorComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.CompressorComplexGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.TurbineComplexGasInput] = core.NewPort()
	result.ports[nodes.TurbineComplexGasInput].SetInnerNode(result)
	result.ports[nodes.TurbineComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.TurbineComplexGasOutput] = core.NewPort()
	result.ports[nodes.TurbineComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.TurbineComplexGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

type TurboCascadeNode interface {
	core.Node
	CompressorComplexGasInput() core.Port
	CompressorComplexGasOutput() core.Port
	TurbineComplexGasInput() core.Port
	TurbineComplexGasOutput() core.Port
	Compressor() constructive.CompressorNode
	Turbine() constructive.TurbineNode
	Transmission() constructive.TransmissionNode
	ProcessCompressor() error
	ProcessTurbine() error
}

type turboCascadeNode struct {
	ports        core.PortsType
	compressor   constructive.CompressorNode
	turbine      constructive.BlockedTurbineNode
	transmission constructive.TransmissionNode
	powerSink    nodes.PowerSink
}

func (node *turboCascadeNode) Compressor() constructive.CompressorNode {
	return node.compressor
}

func (node *turboCascadeNode) Turbine() constructive.TurbineNode {
	return node.turbine
}

func (node *turboCascadeNode) Transmission() constructive.TransmissionNode {
	return node.transmission
}

func (node *turboCascadeNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		CompressorGasInputState  core.PortState                  `json:"compressor_gas_input_state"`
		CompressorGasOutputState core.PortState                  `json:"compressor_gas_output_state"`
		TurbineGasInputState     core.PortState                  `json:"turbine_gas_input_state"`
		TurbineGasOutputState    core.PortState                  `json:"turbine_gas_output_state"`
		Compressor               constructive.CompressorNode     `json:"compressor"`
		CompressorTurbine        constructive.BlockedTurbineNode `json:"compressor_turbine"`
		Transmission             constructive.TransmissionNode   `json:"transmission"`
	}{
		CompressorGasInputState:  node.compressorComplexGasInput().GetState(),
		CompressorGasOutputState: node.compressorComplexGasOutput().GetState(),
		TurbineGasInputState:     node.turbineComplexGasInput().GetState(),
		TurbineGasOutputState:    node.turbineComplexGasOutput().GetState(),
		Compressor:               node.compressor,
		CompressorTurbine:        node.turbine,
		Transmission:             node.transmission,
	})
}

func (node *turboCascadeNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *turboCascadeNode) ProcessCompressor() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	node.writeOutput()

	return nil
}

func (node *turboCascadeNode) ProcessTurbine() error {
	node.readInput()
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) Process() error {
	node.readInput()
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.transmission.Process(); err != nil {
		return err
	}
	if err := node.turbine.Process(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *turboCascadeNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.CompressorComplexGasInput, nodes.TurbineComplexGasInput}, nil
}

func (node *turboCascadeNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.CompressorComplexGasOutput, nodes.TurbineComplexGasOutput}, nil
}

func (node *turboCascadeNode) GetPortTags() []string {
	return []string{
		nodes.CompressorComplexGasInput,
		nodes.CompressorComplexGasOutput,
		nodes.CompressorComplexGasInput,
		nodes.TurbineComplexGasOutput,
	}
}

func (node *turboCascadeNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.CompressorComplexGasInput:
		return node.compressorComplexGasInput(), nil
	case nodes.CompressorComplexGasOutput:
		return node.compressorComplexGasOutput(), nil
	case nodes.TurbineComplexGasInput:
		return node.turbineComplexGasInput(), nil
	case nodes.TurbineComplexGasOutput:
		return node.turbineComplexGasOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *turboCascadeNode) ContextDefined() bool {
	return true
}

func (node *turboCascadeNode) TurbineComplexGasOutput() core.Port {
	return node.turbineComplexGasOutput()
}

func (node *turboCascadeNode) TurbineComplexGasInput() core.Port {
	return node.turbineComplexGasInput()
}

func (node *turboCascadeNode) CompressorComplexGasOutput() core.Port {
	return node.compressorComplexGasOutput()
}

func (node *turboCascadeNode) CompressorComplexGasInput() core.Port {
	return node.compressorComplexGasInput()
}

func (node *turboCascadeNode) linkPorts() {
	core.Link(node.compressor.PowerOutput(), node.transmission.PowerInput())
	core.Link(node.transmission.PowerOutput(), node.turbine.PowerInput())
	core.Link(node.turbine.PowerOutput(), node.powerSink.PowerInput())
}

func (node *turboCascadeNode) readInput() {
	node.compressor.ComplexGasInput().SetState(node.compressorComplexGasInput().GetState())
	node.turbine.ComplexGasInput().SetState(node.turbineComplexGasInput().GetState())
}

func (node *turboCascadeNode) writeOutput() {
	node.compressorComplexGasOutput().SetState(node.compressor.ComplexGasOutput().GetState())
	node.turbineComplexGasOutput().SetState(node.turbine.ComplexGasOutput().GetState())
}

func (node *turboCascadeNode) compressorComplexGasInput() core.Port {
	return node.ports[nodes.CompressorComplexGasInput]
}

func (node *turboCascadeNode) compressorComplexGasOutput() core.Port {
	return node.ports[nodes.CompressorComplexGasOutput]
}

func (node *turboCascadeNode) turbineComplexGasInput() core.Port {
	return node.ports[nodes.TurbineComplexGasInput]
}

func (node *turboCascadeNode) turbineComplexGasOutput() core.Port {
	return node.ports[nodes.TurbineComplexGasOutput]
}
