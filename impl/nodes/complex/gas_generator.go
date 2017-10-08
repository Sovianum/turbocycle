package complex

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/nodes"
	"github.com/Sovianum/turbocycle/impl/nodes/constructive"
	"github.com/Sovianum/turbocycle/fuel"
)

func NewGasGenerator(
	compressorEtaAd, piStag float64,
	fuel fuel.GasFuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0  float64,
	etaT, lambdaOut float64, turbineMassRateRelFunc func(constructive.TurbineNode) float64,

	precision float64,
) GasGenerator {
	var result = &gasGenerator{
		ports:     make(core.PortsType),
		compressor: constructive.NewCompressorNode(compressorEtaAd, piStag, precision),
		burner: constructive.NewBurnerNode(fuel, tgStag, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision),
		compressorTurbine:constructive.NewBlockedTurbineNode(etaT, lambdaOut, precision, turbineMassRateRelFunc),
	}

	result.linkPorts()
	result.ports[nodes.ComplexGasInput] = result.compressor.ComplexGasInput()
	result.ports[nodes.ComplexGasOutput] = result.compressorTurbine.ComplexGasOutput()

	return result
}

type GasGenerator interface {
	core.Node
	nodes.ComplexGasChannel
}

type gasGenerator struct {
	ports             core.PortsType
	compressor        constructive.CompressorNode
	burner            constructive.BurnerNode
	compressorTurbine constructive.BlockedTurbineNode
	powerSink         nodes.PowerSink
}

func (node *gasGenerator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState     core.PortState                  `json:"gas_input_state"`
		GasOutputState    core.PortState                  `json:"gas_output_state"`
		Compressor        constructive.CompressorNode     `json:"compressor"`
		Burner            constructive.BurnerNode         `json:"burner"`
		CompressorTurbine constructive.BlockedTurbineNode `json:"compressor_turbine"`
	}{
		GasInputState:     node.complexGasInput().GetState(),
		GasOutputState:    node.complexGasOutput().GetState(),
		Compressor:        node.compressor,
		Burner:            node.burner,
		CompressorTurbine: node.compressorTurbine,
	})
}

func (node *gasGenerator) GetPorts() core.PortsType {
	return node.ports
}

func (node *gasGenerator) Process() error {
	if err := node.compressor.Process(); err != nil {
		return err
	}
	if err := node.burner.Process(); err != nil {
		return err
	}
	if err := node.compressorTurbine.Process(); err != nil {
		return err
	}
	return nil
}

func (node *gasGenerator) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *gasGenerator) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput}, nil
}

func (node *gasGenerator) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.ComplexGasOutput}
}

func (node *gasGenerator) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.complexGasInput(), nil
	case nodes.ComplexGasOutput:
		return node.complexGasOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *gasGenerator) ContextDefined() bool {
	return true
}

func (node *gasGenerator) ComplexGasInput() core.Port {
	return node.complexGasInput()
}

func (node *gasGenerator) ComplexGasOutput() core.Port {
	return node.complexGasOutput()
}

func (node *gasGenerator) linkPorts() {
	core.Link(node.compressor.ComplexGasOutput(), node.burner.ComplexGasInput())
	core.Link(node.burner.ComplexGasOutput(), node.compressorTurbine.ComplexGasInput())
	core.Link(node.compressor.PowerOutput(), node.compressorTurbine.PowerInput())
	core.Link(node.compressorTurbine.PowerOutput(), node.powerSink.PowerInput())
}

func (node *gasGenerator) complexGasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *gasGenerator) complexGasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}
