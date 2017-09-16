package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

type FreeTurbineNode interface {
	TurbineNode
	LSpecific() float64
}

type freeTurbineNode struct {
	ports           core.PortsType
	etaT            float64
	precision       float64
	lambdaOut       float64
	massRateRelFunc func(TurbineNode) float64
}

func NewFreeTurbineNode(etaT, lambdaOut, precision float64, massRateRelFunc func(TurbineNode) float64) FreeTurbineNode {
	var result = &freeTurbineNode{
		ports:           make(core.PortsType),
		etaT:            etaT,
		precision:       precision,
		lambdaOut:       lambdaOut,
		massRateRelFunc: massRateRelFunc,
	}

	result.ports[complexGasInput] = core.NewPort()
	result.ports[complexGasInput].SetInnerNode(result)
	result.ports[complexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[complexGasOutput] = core.NewPort()
	result.ports[complexGasOutput].SetInnerNode(result)
	result.ports[complexGasOutput].SetState(states.StandardAtmosphereState())

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetInnerNode(result)
	result.ports[powerOutput].SetState(states.StandardPowerState())

	return result
}

func (node *freeTurbineNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState `json:"gas_input_state"`
		GasOutputState   core.PortState `json:"gas_output_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		PiStag           float64        `json:"pi_stag"`
		LSpecific        float64        `json:"l_specific"`
		Eta              float64        `json:"eta"`
	}{
		GasInputState:    node.gasInput().GetState(),
		GasOutputState:   node.gasOutput().GetState(),
		PowerOutputState: node.powerOutput().GetState(),
		PiStag:           node.PiTStag(),
		LSpecific:        node.lSpecific(),
		Eta:              node.etaT,
	})
}

func (node *freeTurbineNode) ContextDefined() bool {
	return true
}

func (node *freeTurbineNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case complexGasInput:
		return node.gasInput(), nil
	case complexGasOutput:
		return node.gasOutput(), nil
	case powerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *freeTurbineNode) GetRequirePortTags() ([]string, error) {
	return []string{complexGasInput, complexGasOutput}, nil
}

func (node *freeTurbineNode) GetUpdatePortTags() ([]string, error) {
	return []string{complexGasOutput, powerOutput}, nil
}

func (node *freeTurbineNode) GetPortTags() []string {
	return []string{complexGasInput, complexGasOutput, powerOutput}
}

func (node *freeTurbineNode) ComplexGasInput() core.Port {
	return node.gasInput()
}

func (node *freeTurbineNode) ComplexGasOutput() core.Port {
	return node.gasOutput()
}

func (node *freeTurbineNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *freeTurbineNode) InputGas() gases.Gas {
	return node.inputGas()
}

func (node *freeTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *freeTurbineNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *freeTurbineNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *freeTurbineNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *freeTurbineNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *freeTurbineNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *freeTurbineNode) PiTStag() float64 {
	return node.piTStag()
}

func (node *freeTurbineNode) LSpecific() float64 {
	return node.lSpecific()
}

func (node *freeTurbineNode) Process() error {
	var gasState = node.gasInput().GetState().(states.ComplexGasPortState)
	gasState.TStag = node.getTStagOut()
	gasState.PStag = node.pStagOut()
	gasState.MassRateRel *= 1 + node.massRateRelFunc(node)

	node.gasOutput().SetState(gasState)

	node.powerOutput().SetState(
		states.NewPowerPortState(node.lSpecific()),
	)

	return nil
}

func (node *freeTurbineNode) lSpecific() float64 {
	return gases.CpMean(node.inputGas(), node.tStagIn(), node.tStagOut(), defaultN) * (node.tStagIn() - node.tStagOut())
}

func (node *freeTurbineNode) getTStagOut() float64 {
	var tStagOutCurr = node.tStagOutNext(
		node.pStagIn(), node.pStagOut(), node.tStagIn(), node.tStagIn(),
	)
	var tStagOutNext = node.tStagOutNext(
		node.pStagIn(), node.pStagOut(), node.tStagIn(), tStagOutCurr,
	)

	for !common.Converged(tStagOutCurr, tStagOutNext, node.precision) {
		tStagOutCurr = tStagOutNext
		node.tStagOutNext(
			node.pStagIn(), node.pStagOut(), node.tStagIn(), tStagOutCurr,
		)
	}

	return tStagOutNext
}

func (node *freeTurbineNode) tStagOutNext(pStagIn, pStagOut, tStagIn, tStagOutCurr float64) float64 {
	var k = gases.KMean(node.inputGas(), tStagIn, tStagOutCurr, defaultN)
	var piT = pStagIn / pStagOut
	var x = math.Pow(piT, (1-k)/k)

	return tStagIn * (1 - (1-x)*node.etaT)
}

func (node *freeTurbineNode) piTStag() float64 {
	return node.pStagIn() / node.pStagOut()
}

func (node *freeTurbineNode) inputGas() gases.Gas {
	return node.gasInput().GetState().(states.ComplexGasPortState).Gas
}

func (node *freeTurbineNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *freeTurbineNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.ComplexGasPortState).PStag
}

func (node *freeTurbineNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).TStag
}

func (node *freeTurbineNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).PStag
}

func (node *freeTurbineNode) gasInput() core.Port {
	return node.ports[complexGasInput]
}

func (node *freeTurbineNode) gasOutput() core.Port {
	return node.ports[complexGasOutput]
}

func (node *freeTurbineNode) powerOutput() core.Port {
	return node.ports[powerOutput]
}
