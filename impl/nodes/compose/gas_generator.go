package compose

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/states"
)

func NewGasGeneratorNode(
	compressorEtaAd, piStag float64,
	fuel fuel.GasFuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0 float64,
	etaT, lambdaOut float64, turbineMassRateRelFunc func(constructive.TurbineNode) float64,
	etaM float64,
	precision float64,
) GasGeneratorNode {
	var result = &gasGeneratorNode{
		ports:        make(core.PortsType),
		turboCascade: NewTurboCascadeNode(compressorEtaAd, piStag, etaT, lambdaOut, turbineMassRateRelFunc, etaM, precision),
		burner:       constructive.NewBurnerNode(fuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision),
	}

	result.linkPorts()

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

type GasGeneratorNode interface {
	core.Node
	nodes.ComplexGasChannel
}

type gasGeneratorNode struct {
	ports        core.PortsType
	burner       constructive.BurnerNode
	turboCascade TurboCascadeNode
}

func (node *gasGeneratorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState  core.PortState          `json:"gas_input_state"`
		GasOutputState core.PortState          `json:"gas_output_state"`
		Burner         constructive.BurnerNode `json:"burner"`
		TurboCascade   TurboCascadeNode        `json:"turbo_cascade"`
	}{
		GasInputState:  node.complexGasInput().GetState(),
		GasOutputState: node.complexGasOutput().GetState(),
		Burner:         node.burner,
		TurboCascade:   node.turboCascade,
	})
}

func (node *gasGeneratorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasGeneratorNode) Process() error {
	node.readInput()
	if err := node.turboCascade.ProcessCompressor(); err != nil {
		return err
	}
	if err := node.burner.Process(); err != nil {
		return err
	}
	if err := node.turboCascade.ProcessTurbine(); err != nil {
		return err
	}
	node.writeOutput()
	return nil
}

func (node *gasGeneratorNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *gasGeneratorNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput}, nil
}

func (node *gasGeneratorNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.ComplexGasOutput}
}

func (node *gasGeneratorNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.ComplexGasOutput:
		return node.complexGasOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *gasGeneratorNode) ContextDefined() bool {
	return true
}

func (node *gasGeneratorNode) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *gasGeneratorNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput()
}

func (node *gasGeneratorNode) readInput() {
	node.turboCascade.CompressorComplexGasInput().SetState(node.complexGasInput().GetState())
}

func (node *gasGeneratorNode) writeOutput() {
	node.complexGasOutput().SetState(node.turboCascade.TurbineComplexGasOutput().GetState())
}

func (node *gasGeneratorNode) linkPorts() {
	core.Link(node.turboCascade.CompressorComplexGasOutput(), node.burner.ComplexGasInput())
	core.Link(node.burner.ComplexGasOutput(), node.turboCascade.TurbineComplexGasInput())
}

func (node *gasGeneratorNode) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *gasGeneratorNode) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}
