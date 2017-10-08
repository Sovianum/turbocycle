package compose

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
)

func NewTurboCascadeNode(
	compressorEtaAd, piStag float64,
	etaT, lambdaOut float64, turbineMassRateRelFunc func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) TurboCascadeNode {
	var result = &turboCascadeNode{
		ports:             make(core.PortsType),
		compressor:        constructive.NewCompressorNode(compressorEtaAd, piStag, precision),
		compressorTurbine: constructive.NewBlockedTurbineNode(etaT, lambdaOut, precision, turbineMassRateRelFunc),
		transmission:      constructive.NewTransmissionNode(etaM),
	}

	result.linkPorts()
	result.ports[nodes.ComplexGasInput] = result.compressor.ComplexGasInput()
	result.ports[nodes.ComplexGasOutput] = result.compressorTurbine.ComplexGasOutput()

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
}

type turboCascadeNode struct {
	ports             core.PortsType
	compressor        constructive.CompressorNode
	compressorTurbine constructive.BlockedTurbineNode
	transmission      constructive.TransmissionNode
	powerSink         nodes.PowerSink
}

func (node *turboCascadeNode) Compressor() constructive.CompressorNode {
	return node.compressor
}

func (node *turboCascadeNode) Turbine() constructive.TurbineNode {
	return node.compressorTurbine
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
		CompressorTurbine:        node.compressorTurbine,
		Transmission:             node.transmission,
	})
}

func (node *turboCascadeNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *turboCascadeNode) Process() error {
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.compressorTurbine.Process(); err != nil {
		return err
	}
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
	core.Link(node.transmission.PowerOutput(), node.compressorTurbine.PowerInput())
	core.Link(node.compressorTurbine.PowerOutput(), node.powerSink.PowerInput())
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
	return node.ports[nodes.TurbineComplexGasInput]
}
