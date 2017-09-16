package constructive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
	"github.com/Sovianum/turbocycle/impl/nodes"
)

type CompressorNode interface {
	core.Node
	nodes.ComplexGasChannel
	nodes.PowerSource
	LSpecific() float64
}

// TODO add collector port
type compressorNode struct {
	ports     core.PortsType
	EtaAd     float64
	Precision float64
	PiStag    float64
}

func (node *compressorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState `json:"gas_input_state"`
		GasOutputState   core.PortState `json:"gas_output_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		EtaAd            float64        `json:"eta_ad"`
		PiStag           float64        `json:"pi_stag"`
		MassRateRel      float64        `json:"mass_rate_rel"`
	}{
		GasInputState:    node.gasInput().GetState(),
		GasOutputState:   node.gasOutput().GetState(),
		PowerOutputState: node.powerOutput().GetState(),
		EtaAd:            node.EtaAd,
		PiStag:           node.PiStag,
		MassRateRel:      node.gasInput().GetState().(states.ComplexGasPortState).MassRateRel,
	})
}

func NewCompressorNode(etaAd, piStag, precision float64) CompressorNode {
	var result = &compressorNode{
		ports:     make(core.PortsType),
		EtaAd:     etaAd,
		Precision: precision,
		PiStag:    piStag,
	}

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.PowerOutput] = core.NewPort()
	result.ports[nodes.PowerOutput].SetInnerNode(result)
	result.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	return result
}

func (node *compressorNode) ContextDefined() bool {
	return true
}

func (node *compressorNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case nodes.ComplexGasInput:
		return node.gasInput(), nil
	case nodes.ComplexGasOutput:
		return node.gasOutput(), nil
	case nodes.PowerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *compressorNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput}, nil
}

func (node *compressorNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput, nodes.PowerOutput}, nil
}

func (node *compressorNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.ComplexGasOutput, nodes.PowerOutput}
}

func (node *compressorNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *compressorNode) Process() error {
	if node.PiStag <= 1 {
		return errors.New(fmt.Sprintf("Invalid piStag = %f", node.PiStag))
	}

	var pStagOut = node.pStagIn() * node.PiStag
	var tStagOut = node.getTStagOut(node.PiStag, node.tStagIn(), node.tStagIn())

	var gasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)
	gasState.TStag = tStagOut
	gasState.PStag = pStagOut

	node.gasOutput().SetState(gasState)

	node.powerOutput().SetState(states.NewPowerPortState(-node.lSpecific()))
	// TODO add and set collector port

	return nil
}

func (node *compressorNode) ComplexGasInput() core.Port {
	return node.gasInput()
}

func (node *compressorNode) ComplexGasOutput() core.Port {
	return node.gasOutput()
}

func (node *compressorNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *compressorNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *compressorNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *compressorNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *compressorNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *compressorNode) LSpecific() float64 {
	return node.lSpecific()
}

func (node *compressorNode) lSpecific() float64 {
	var cpMean = gases.CpMean(node.gas(), node.tStagIn(), node.tStagOut(), nodes.DefaultN)
	return cpMean * (node.tStagOut() - node.tStagIn())
}

func (node *compressorNode) getTStagOut(piCStag, tStagIn, tStagOutInit float64) float64 {
	var k = gases.K(node.gas(), tStagIn)
	var x = math.Pow(piCStag, (k-1)/k)

	var tOutCurr = tStagIn * (1 + (x-1)/node.EtaAd)
	var tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)

	for !common.Converged(tOutCurr, tOutNext, node.Precision) {
		tOutCurr = tOutNext
		tOutNext = node.tStagOutNewFunc(piCStag, tStagIn, tStagOutInit)
	}

	return tOutNext
}

func (node *compressorNode) tStagOutNewFunc(piCStag, tStagIn, tStagOutCurr float64) float64 {
	var x = node.xFunc(piCStag, tStagIn, tStagOutCurr)
	return tStagIn * (1 + (x-1)/node.EtaAd)
}

func (node *compressorNode) xFunc(piCStag, tStagIn, tStagOut float64) float64 {
	var k = gases.KMean(node.gas(), tStagIn, tStagOut, nodes.DefaultN)
	return math.Pow(piCStag, (k-1)/k)
}

func (node *compressorNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *compressorNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).TStag
}

func (node *compressorNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.ComplexGasPortState).PStag
}

func (node *compressorNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).PStag
}

func (node *compressorNode) gas() gases.Gas {
	return node.ports[nodes.ComplexGasInput].GetState().(states.ComplexGasPortState).Gas
}

func (node *compressorNode) gasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *compressorNode) gasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *compressorNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
