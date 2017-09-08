package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

type freeTurbineNode struct {
	ports           core.PortsType
	etaT            float64
	precision       float64
	lambdaOut       float64
	massRateRelFunc func(TurbineNode) float64
}

func NewFreeTurbineNode(etaT, lambdaOut, precision float64, massRateRelFunc func(TurbineNode) float64) *freeTurbineNode {
	var result = &freeTurbineNode{
		ports:           make(core.PortsType),
		etaT:            etaT,
		precision:       precision,
		lambdaOut:       lambdaOut,
		massRateRelFunc: massRateRelFunc,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetInnerNode(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetInnerNode(result)

	return result
}

func (node *freeTurbineNode) GetPortByTag(tag string) (*core.Port, error) {
	switch tag {
	case gasInput:
		return node.gasInput(), nil
	case gasOutput:
		return node.gasOutput(), nil
	case powerOutput:
		return node.PowerOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("Port with tag \"%s\" not found", tag))
	}
}

func (node *freeTurbineNode) GetRequirePortTags() []string {
	return []string{gasInput, gasOutput}
}

func (node *freeTurbineNode) GetUpdatePortTags() []string {
	return []string{gasOutput, powerOutput}
}

func (node *freeTurbineNode) GetPortTags() []string {
	return []string{gasInput, gasOutput, powerOutput}
}

func (node *freeTurbineNode) GasInput() *core.Port {
	return node.gasInput()
}

func (node *freeTurbineNode) GasOutput() *core.Port {
	return node.gasOutput()
}

func (node *freeTurbineNode) PowerOutput() *core.Port {
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

func (node *freeTurbineNode) Pit() float64 {
	return node.pit()
}

func (node *freeTurbineNode) Process() error {
	var gasState = node.gasInput().GetState().(states.GasPortState)
	gasState.TStag = node.getTStagOut()
	gasState.PStag = node.pStagOut()
	gasState.MassRateRel *= 1 + node.massRateRelFunc(node)

	node.gasOutput().SetState(gasState)

	node.powerOutput().SetState(
		states.NewPowerPortState(node.turbineLabour()),
	)

	return nil
}

func (node *freeTurbineNode) turbineLabour() float64 {
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

func (node *freeTurbineNode) pit() float64 {
	return node.pStagIn() / node.pStagOut()
}

func (node *freeTurbineNode) inputGas() gases.Gas {
	return node.gasInput().GetState().(states.GasPortState).Gas
}

func (node *freeTurbineNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *freeTurbineNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}

func (node *freeTurbineNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *freeTurbineNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *freeTurbineNode) gasInput() *core.Port {
	return node.ports[gasInput]
}

func (node *freeTurbineNode) gasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *freeTurbineNode) powerOutput() *core.Port {
	return node.ports[powerOutput]
}
