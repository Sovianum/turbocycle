package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/gdf"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

type BlockedTurbineNode interface {
	TurbineNode
	PowerSink
}

type blockedTurbineNode struct {
	ports           core.PortsType
	etaT            float64
	precision       float64
	lambdaOut       float64
	massRateRelFunc func(TurbineNode) float64
}

func NewBlockedTurbineNode(etaT, lambdaOut, precision float64, massRateRelFunc func(TurbineNode) float64) BlockedTurbineNode {
	var result = &blockedTurbineNode{
		ports:           make(core.PortsType),
		etaT:            etaT,
		precision:       precision,
		lambdaOut:       lambdaOut,
		massRateRelFunc: massRateRelFunc,
	}

	result.ports[powerInput] = core.NewPort()
	result.ports[powerInput].SetInnerNode(result)
	result.ports[powerInput].SetState(states.StandardPowerState())

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetInnerNode(result)
	result.ports[powerOutput].SetState(states.StandardPowerState())

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetInnerNode(result)
	result.ports[gasInput].SetState(states.StandardAtmosphereState())

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)
	result.ports[gasOutput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *blockedTurbineNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState `json:"gas_input_state"`
		GasOutputState   core.PortState `json:"gas_output_state"`
		PowerInputState  core.PortState `json:"power_input_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		LSpecific        float64         `json:"l_specific"`
		Eta              float64         `json:"eta"`
	}{
		GasInputState:    node.gasInput().GetState(),
		GasOutputState:   node.gasOutput().GetState(),
		PowerInputState:  node.powerInput().GetState(),
		PowerOutputState: node.powerOutput().GetState(),
		LSpecific:        node.turbineLabour(),
		Eta:              node.etaT,
	})
}

func (node *blockedTurbineNode) ContextDefined() bool {
	return true
}

func (node *blockedTurbineNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case gasInput:
		return node.gasInput(), nil
	case gasOutput:
		return node.gasOutput(), nil
	case powerInput:
		return node.powerInput(), nil
	case powerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *blockedTurbineNode) GetRequirePortTags() ([]string, error) {
	return []string{gasInput, powerInput}, nil
}

func (node *blockedTurbineNode) GetUpdatePortTags() ([]string, error) {
	return []string{gasOutput, powerOutput}, nil
}

func (node *blockedTurbineNode) GetPortTags() []string {
	return []string{gasInput, powerInput, gasOutput, powerOutput}
}

func (node *blockedTurbineNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *blockedTurbineNode) Process() error {
	var gasState = node.GasInput().GetState().(states.GasPortState)
	gasState.TStag = node.getTStagOut(node.turbineLabour())

	var piTStag = node.piTStag(gasState.TStag)
	var pi = gdf.Pi(node.lambdaOut, gases.KMean(node.inputGas(), node.tStagIn(), gasState.TStag, defaultN))
	gasState.PStag = node.pStagIn() / (piTStag * pi)
	gasState.MassRateRel *= 1 + node.massRateRelFunc(node)

	node.gasOutput().SetState(gasState)
	node.powerOutput().SetState(states.NewPowerPortState(node.turbineLabour())) // TODO maybe need to pass sum of labours

	return nil
}

func (node *blockedTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *blockedTurbineNode) InputGas() gases.Gas {
	return node.inputGas()
}

func (node *blockedTurbineNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *blockedTurbineNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *blockedTurbineNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *blockedTurbineNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *blockedTurbineNode) PiTStag() float64 {
	return node.piTStag(node.tStagOut())
}

func (node *blockedTurbineNode) GasInput() core.Port {
	return node.gasInput()
}

func (node *blockedTurbineNode) GasOutput() core.Port {
	return node.gasOutput()
}

func (node *blockedTurbineNode) PowerInput() core.Port {
	return node.powerInput()
}

func (node *blockedTurbineNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *blockedTurbineNode) getTStagOut(turbineLabour float64) float64 {
	var tTStagCurr = node.getInitTtStag(node.turbineLabour())
	var tTStagNew = node.getNewTtStag(tTStagCurr, node.turbineLabour())

	for !common.Converged(tTStagCurr, tTStagNew, node.precision) {
		tTStagCurr = tTStagNew
		tTStagNew = node.getNewTtStag(tTStagCurr, node.turbineLabour())
	}

	return tTStagNew
}

func (node *blockedTurbineNode) getInitTtStag(turbineLabour float64) float64 {
	return node.getNewTtStag(0.8*node.tStagIn(), turbineLabour) // TODO move 0.8 out of code
}

func (node *blockedTurbineNode) getNewTtStag(currTtStag, turbineLabour float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), currTtStag, defaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), currTtStag, defaultN)

	var piTStag = node.getPiTStag(k, cp, turbineLabour)

	return node.tStagIn() * (1 - (1-math.Pow(piTStag, (1-k)/k))*node.etaT)
}

func (node *blockedTurbineNode) inputGas() gases.Gas {
	return node.gasInput().GetState().(states.GasPortState).Gas
}

func (node *blockedTurbineNode) piTStag(tStagOut float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), tStagOut, defaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), tStagOut, defaultN)

	return node.getPiTStag(k, cp, node.turbineLabour())
}

func (node *blockedTurbineNode) getPiTStag(k, cp, turbineLabour float64) float64 {
	return math.Pow(
		1-turbineLabour/(cp*node.tStagIn()*node.etaT),
		k/(1-k),
	)
}

func (node *blockedTurbineNode) turbineLabour() float64 {
	return -node.powerInput().GetState().(states.PowerPortState).LSpecific
}

func (node *blockedTurbineNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *blockedTurbineNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}

func (node *blockedTurbineNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *blockedTurbineNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *blockedTurbineNode) gasInput() core.Port {
	return node.ports[gasInput]
}

func (node *blockedTurbineNode) gasOutput() core.Port {
	return node.ports[gasOutput]
}

func (node *blockedTurbineNode) powerInput() core.Port {
	return node.ports[powerInput]
}

func (node *blockedTurbineNode) powerOutput() core.Port {
	return node.ports[powerOutput]
}
