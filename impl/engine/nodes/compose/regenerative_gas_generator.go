package compose

import (
	"encoding/json"
	"fmt"

	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/fuel"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

func NewRegenerativeGasGeneratorNode(
	compressorEtaAd, piStag float64,
	fuel fuel.GasFuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0 float64,
	etaT, lambdaOut float64, turbineMassRateRelFunc func(constructive.TurbineNode) float64,
	sigmaRegenerator float64,
	sigmaRegeneratorPipe float64,
	etaM float64,
	precision float64,
) RegenerativeGasGeneratorNode {
	var result = &regenerativeGasGeneratorNode{
		ports:make(core.PortsType),
		turboCascade:NewTurboCascadeNode(compressorEtaAd, piStag, etaT, lambdaOut, turbineMassRateRelFunc, etaM, precision),
		burner:constructive.NewBurnerNode(fuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision),
		regenerator:constructive.NewRegeneratorNode(sigmaRegenerator, precision, constructive.SigmaByColdSide),
		regeneratorPressureDrop:constructive.NewPressureLossNode(sigmaRegeneratorPipe),
	}

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.HeatExchangerHotInput] = core.NewPort()
	result.ports[nodes.HeatExchangerHotInput].SetInnerNode(result)
	result.ports[nodes.HeatExchangerHotInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.HeatExchangerHotOutput] = core.NewPort()
	result.ports[nodes.HeatExchangerHotOutput].SetInnerNode(result)
	result.ports[nodes.HeatExchangerHotOutput].SetState(states.StandardAtmosphereState())

	result.linkPorts()
	return result
}

type RegenerativeGasGeneratorNode interface {
	core.Node
	nodes.ComplexGasChannel
	HeatExchangerHotInput() core.Port
	HeatExchangerHotOutput() core.Port
	Burner() constructive.BurnerNode
	TurboCascade() TurboCascadeNode
	Regenerator() constructive.RegeneratorNode
}

type regenerativeGasGeneratorNode struct {
	ports                   core.PortsType
	burner                  constructive.BurnerNode
	turboCascade            TurboCascadeNode
	regenerator             constructive.RegeneratorNode
	regeneratorPressureDrop constructive.PressureLossNode
}

func (node *regenerativeGasGeneratorNode) Burner() constructive.BurnerNode {
	return node.burner
}

func (node *regenerativeGasGeneratorNode) TurboCascade() TurboCascadeNode {
	return node.turboCascade
}

func (node *regenerativeGasGeneratorNode) Regenerator() constructive.RegeneratorNode {
	return node.regenerator
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
	if err := node.burner.Process(); err != nil {
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
