package nodes

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/gdf"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

type blockedTurbineNode struct {
	ports           core.PortsType
	etaT            float64
	precision       float64
	lambdaOut       float64
	massRateRelFunc func(TurbineNode) float64
}

func NewBlockedTurbineNode(etaT, lambdaOut, precision float64, massRateRelFunc func(TurbineNode) float64) *blockedTurbineNode {
	var result = &blockedTurbineNode{
		ports:           make(core.PortsType),
		etaT:            etaT,
		precision:       precision,
		lambdaOut:       lambdaOut,
		massRateRelFunc: massRateRelFunc,
	}

	result.ports[powerInput] = core.NewPort()
	result.ports[powerInput].SetDest(result)

	result.ports[powerOutput] = core.NewPort()
	result.ports[powerOutput].SetSrc(result)

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetDest(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetSrc(result)

	return result
}

func NewBlockedTurbineNodeShort(etaT float64, massRateRel func(TurbineNode) float64) *blockedTurbineNode {
	return NewBlockedTurbineNode(etaT, 0.3, 0.05, massRateRel) // TODO remove hardcoded constants
}

func (node *blockedTurbineNode) GetPorts() core.PortsType {
	return node.ports
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

func (node *blockedTurbineNode) Pit() float64 {
	return node.pit(node.tStagOut())
}

func (node *blockedTurbineNode) GasInput() *core.Port {
	return node.gasInput()
}

func (node *blockedTurbineNode) GasOutput() *core.Port {
	return node.gasOutput()
}

func (node *blockedTurbineNode) PowerInput() *core.Port {
	return node.powerInput()
}

func (node *blockedTurbineNode) PowerOutput() *core.Port {
	return node.powerOutput()
}

func (node *blockedTurbineNode) Process() error {
	var gasState = node.GasInput().GetState().(states.GasPortState)
	gasState.TStag = node.getTStagOut(node.turbineLabour())

	var pit = node.pit(gasState.TStag)
	var pi = gdf.Pi(node.lambdaOut, gases.KMean(node.inputGas(), node.tStagIn(), gasState.TStag, defaultN))
	gasState.PStag = node.pStagIn() / (pit * pi)
	gasState.MassRateRel *= 1 + node.massRateRelFunc(node)

	node.gasOutput().SetState(gasState)
	node.powerOutput().SetState(states.NewPowerPortState(node.turbineLabour())) // TODO maybe need to pass sum of labours

	return nil
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

	var pit = node.getPit(k, cp, turbineLabour)

	return node.tStagIn() * (1 - (1-math.Pow(pit, (1-k)/k))*node.etaT)
}

func (node *blockedTurbineNode) getPit(k, cp, turbineLabour float64) float64 {
	return math.Pow(
		1-turbineLabour/(cp*node.tStagIn()*node.etaT),
		k/(1-k),
	)
}

func (node *blockedTurbineNode) inputGas() gases.Gas {
	return node.gasInput().GetState().(states.GasPortState).Gas
}

func (node *blockedTurbineNode) pit(tStagOut float64) float64 {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), tStagOut, defaultN)
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), tStagOut, defaultN)

	return math.Pow(
		1-node.turbineLabour()/(cp*node.tStagIn()*node.etaT),
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

func (node *blockedTurbineNode) gasInput() *core.Port {
	return node.ports[gasInput]
}

func (node *blockedTurbineNode) gasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *blockedTurbineNode) powerInput() *core.Port {
	return node.ports[powerInput]
}

func (node *blockedTurbineNode) powerOutput() *core.Port {
	return node.ports[powerOutput]
}
