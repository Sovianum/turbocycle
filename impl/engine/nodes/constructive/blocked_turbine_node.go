package constructive

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/helpers/gdf"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
)

type BlockedTurbineNode interface {
	TurbineNode
	nodes.ComplexGasSource
	nodes.PowerSink
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

	result.ports[nodes.PowerInput] = core.NewPort()
	result.ports[nodes.PowerInput].SetInnerNode(result)
	result.ports[nodes.PowerInput].SetState(states.StandardPowerState())

	result.ports[nodes.PowerOutput] = core.NewPort()
	result.ports[nodes.PowerOutput].SetInnerNode(result)
	result.ports[nodes.PowerOutput].SetState(states.StandardPowerState())

	result.ports[nodes.ComplexGasInput] = core.NewPort()
	result.ports[nodes.ComplexGasInput].SetInnerNode(result)
	result.ports[nodes.ComplexGasInput].SetState(states.StandardAtmosphereState())

	result.ports[nodes.ComplexGasOutput] = core.NewPort()
	result.ports[nodes.ComplexGasOutput].SetInnerNode(result)
	result.ports[nodes.ComplexGasOutput].SetState(states.StandardAtmosphereState())

	return result
}

func (node *blockedTurbineNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GasInputState    core.PortState `json:"gas_input_state"`
		GasOutputState   core.PortState `json:"gas_output_state"`
		PowerInputState  core.PortState `json:"power_input_state"`
		PowerOutputState core.PortState `json:"power_output_state"`
		LSpecific        float64        `json:"l_specific"`
		Eta              float64        `json:"eta"`
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
	case nodes.ComplexGasInput:
		return node.gasInput(), nil
	case nodes.ComplexGasOutput:
		return node.gasOutput(), nil
	case nodes.PowerInput:
		return node.powerInput(), nil
	case nodes.PowerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, fmt.Errorf("port with tag \"%s\" not found", tag)
	}
}

func (node *blockedTurbineNode) GetRequirePortTags() ([]string, error) {
	return []string{nodes.ComplexGasInput, nodes.PowerInput}, nil
}

func (node *blockedTurbineNode) GetUpdatePortTags() ([]string, error) {
	return []string{nodes.ComplexGasOutput, nodes.PowerOutput}, nil
}

func (node *blockedTurbineNode) GetPortTags() []string {
	return []string{nodes.ComplexGasInput, nodes.PowerInput, nodes.ComplexGasOutput, nodes.PowerOutput}
}

func (node *blockedTurbineNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *blockedTurbineNode) Process() error {
	var gasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)

	var err error
	gasState.TStag, err = node.getTStagOut(node.turbineLabour())
	if err != nil {
		return err
	}

	var piTStag = node.piTStag(gasState.TStag)
	var pi = gdf.Pi(node.lambdaOut, gases.KMean(node.inputGas(), node.tStagIn(), gasState.TStag, nodes.DefaultN))
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

func (node *blockedTurbineNode) ComplexGasInput() core.Port {
	return node.gasInput()
}

func (node *blockedTurbineNode) ComplexGasOutput() core.Port {
	return node.gasOutput()
}

func (node *blockedTurbineNode) PowerInput() core.Port {
	return node.powerInput()
}

func (node *blockedTurbineNode) PowerOutput() core.Port {
	return node.powerOutput()
}

func (node *blockedTurbineNode) getTStagOut(turbineLabour float64) (float64, error) {
	var tTStagCurr = node.getInitTtStag(node.turbineLabour())
	var tTStagNew = node.getNewTtStag(tTStagCurr, node.turbineLabour())

	for !common.Converged(tTStagCurr, tTStagNew, node.precision) {
		if math.IsNaN(tTStagCurr) || math.IsNaN(tTStagNew) {
			return 0, errors.New("failed to converge: try different initial guess")
		}
		tTStagCurr = tTStagNew
		tTStagNew = node.getNewTtStag(tTStagCurr, node.turbineLabour())
	}

	return tTStagNew, nil
}

func (node *blockedTurbineNode) getInitTtStag(turbineLabour float64) float64 {
	return node.getNewTtStag(0.8*node.tStagIn(), turbineLabour) // TODO move 0.8 out of code
}

func (node *blockedTurbineNode) getNewTtStag(currTtStag, turbineLabour float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), currTtStag, nodes.DefaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), currTtStag, nodes.DefaultN)

	var piTStag = node.getPiTStag(k, cp, turbineLabour)

	return node.tStagIn() * (1 - (1-math.Pow(piTStag, (1-k)/k))*node.etaT)
}

func (node *blockedTurbineNode) inputGas() gases.Gas {
	return node.gasInput().GetState().(states.ComplexGasPortState).Gas
}

func (node *blockedTurbineNode) piTStag(tStagOut float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN)

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
	return node.gasInput().GetState().(states.ComplexGasPortState).TStag
}

func (node *blockedTurbineNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.ComplexGasPortState).PStag
}

func (node *blockedTurbineNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).TStag
}

func (node *blockedTurbineNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.ComplexGasPortState).PStag
}

func (node *blockedTurbineNode) gasInput() core.Port {
	return node.ports[nodes.ComplexGasInput]
}

func (node *blockedTurbineNode) gasOutput() core.Port {
	return node.ports[nodes.ComplexGasOutput]
}

func (node *blockedTurbineNode) powerInput() core.Port {
	return node.ports[nodes.PowerInput]
}

func (node *blockedTurbineNode) powerOutput() core.Port {
	return node.ports[nodes.PowerOutput]
}
