package compose

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
)

type RegenerativeGasGeneratorNode interface {
	core.Node
	nodes.ComplexGasChannel
	HeatExchangerHotInput() core.Port
	HeatExchangerHotOutput() core.Port
}

type regenerativeGasGeneratorNode struct {
	ports                   core.PortsType
	burner                  constructive.BurnerNode
	turboCascade            TurboCascadeNode
	regenerator             constructive.RegeneratorNode
	regeneratorPressureDrop constructive.PressureLossNode
}

func (node *regenerativeGasGeneratorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState               core.PortState                `json:"gas_input_state"`
		GasOutputState              core.PortState                `json:"gas_output_state"`
		HeatExchangerHotInputState  core.PortState                `json:"heat_exchanger_hot_input_state"`
		HeatExchangerHotOutputState core.PortState                `json:"heat_exchanger_hot_output_state"`
		Burner                      constructive.BurnerNode       `json:"burner"`
		TurboCascade                TurboCascadeNode              `json:"turbo_cascade"`
		Regenerator                 constructive.RegeneratorNode  `json:"regenerator"`
		RegeneratorPipe             constructive.PressureLossNode `json:"regenerator_pipe"`
	}{
		GasInputState:               node.complexGasInput().GetState(),
		GasOutputState:              node.complexGasOutput().GetState(),
		HeatExchangerHotInputState:  node.heatExchangerHotInput().GetState(),
		HeatExchangerHotOutputState: node.heatExchangerHotOutput().GetState(),
		Burner:          node.burner,
		TurboCascade:    node.turboCascade,
		Regenerator:     node.regenerator,
		RegeneratorPipe: node.regeneratorPressureDrop,
	})
}

func (node *regenerativeGasGeneratorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *regenerativeGasGeneratorNode) Process() error {
	node.readInput()
	if err := node.turboCascade.ProcessCompressor(); err != nil {
		return err
	}
	if err := node.regeneratorPressureDrop.Process(); err != nil {
		return err
	}
	if err := node.regenerator.Process(); err != nil {
		return err
	}
	if err := node.turboCascade.ProcessTurbine(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *regenerativeGasGeneratorNode) GetRequirePortTags() ([]string, error) {
	return []string{
		nodes.ComplexGasInput,
		nodes.HeatExchangerHotInput,
	}, nil
}

func (node *regenerativeGasGeneratorNode) GetUpdatePortTags() ([]string, error) {
	return []string{
		nodes.ComplexGasOutput,
		nodes.HeatExchangerHotOutput,
	}, nil
}

func (node *regenerativeGasGeneratorNode) GetPortTags() []string {
	return []string{
		nodes.ComplexGasInput,
		nodes.ComplexGasOutput,
		nodes.HeatExchangerHotInput,
		nodes.HeatExchangerHotOutput,
	}
}

func (node *regenerativeGasGeneratorNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.ComplexGasOutput:
		return node.complexGasOutput(), nil
	case nodes.HeatExchangerHotInput:
		return node.heatExchangerHotInput(), nil
	case nodes.HeatExchangerHotOutput:
		return node.heatExchangerHotOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *regenerativeGasGeneratorNode) ContextDefined() bool {
	return true
}

func (node *regenerativeGasGeneratorNode) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *regenerativeGasGeneratorNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput()
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotInput() core.Port {
	return node.heatExchangerHotInput()
}

func (node *regenerativeGasGeneratorNode) HeatExchangerHotOutput() core.Port {
	return node.heatExchangerHotOutput()
}

func (node *regenerativeGasGeneratorNode) readInput() {
	node.turboCascade.CompressorComplexGasInput().SetState(node.ComplexGasInput().GetState())
	node.regenerator.HotInput().SetState(node.heatExchangerHotInput().GetState())
}

func (node *regenerativeGasGeneratorNode) writeOutput() {
	node.complexGasOutput().SetState(node.turboCascade.TurbineComplexGasOutput().GetState())
	node.heatExchangerHotOutput().SetState(node.regenerator.HotOutput().GetState())
}

func (node *regenerativeGasGeneratorNode) linkPorts() {
	core.Link(node.turboCascade.CompressorComplexGasOutput(), node.regeneratorPressureDrop.ComplexGasInput())
	core.Link(node.regeneratorPressureDrop.ComplexGasOutput(), node.regenerator.ColdInput())
	core.Link(node.regenerator.ColdOutput(), node.burner.ComplexGasInput())
	core.Link(node.burner.ComplexGasOutput(), node.turboCascade.TurbineComplexGasInput())
}

func (node *regenerativeGasGeneratorNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *regenerativeGasGeneratorNode) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *regenerativeGasGeneratorNode) heatExchangerHotInput() core.Port {
	return node.ports[nodes.HeatExchangerHotInput]
}

func (node *regenerativeGasGeneratorNode) heatExchangerHotOutput() core.Port {
	return node.ports[nodes.HeatExchangerHotOutput]
}
