package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"math"
)

type BurnerNode struct {
	ports     core.PortsType
	fuel      fuel.GasFuel
	outletGas gases.Gas
	tgStag    float64
	tFuel     float64
	sigma     float64
	etaBurn   float64
	initAlpha float64
	t0        float64
	precision float64
	alpha     float64
}

func NewBurnerNode(
	fuel fuel.GasFuel, tgStag, tFuel, sigma, etaBurn, initAlpha, t0, precision float64,
) *BurnerNode {
	var result = &BurnerNode{
		ports:     make(core.PortsType),
		fuel:      fuel,
		tgStag:    tgStag,
		tFuel:     tFuel,
		sigma:     sigma,
		etaBurn:   etaBurn,
		initAlpha: initAlpha,
		t0:        t0,
		precision: precision,
	}

	result.ports[gasInput] = core.NewPort()
	result.ports[gasInput].SetDest(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetSrc(result)

	return result
}

func NewBurnerNodeShort(fuel fuel.GasFuel, tgStag, tFuel, sigma, etaBurn float64) *BurnerNode {
	return NewBurnerNode(
		fuel, tgStag, tFuel, sigma, etaBurn, 3, 290, 0.01,
	)
}

func (node *BurnerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *BurnerNode) GasInput() *core.Port {
	return node.ports[gasInput]
}

func (node *BurnerNode) GasOutput() *core.Port {
	return node.ports[gasOutput]
}

func (node *BurnerNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *BurnerNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *BurnerNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *BurnerNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *BurnerNode) Process() error {
	var fuelMassRateRel, alpha = node.getFuelParameters(node.initAlpha)
	node.alpha = alpha

	var gasState = node.GasInput().GetState().(states.GasPortState)
	gasState.Gas = node.outletGas
	gasState.TStag = node.tgStag
	gasState.PStag = node.pStagIn() * node.sigma
	gasState.MassRateRel *= 1 + fuelMassRateRel

	node.GasOutput().SetState(gasState)

	return nil
}

func (node *BurnerNode) inletGas() gases.Gas {
	return node.GasInput().GetState().(states.GasPortState).Gas
}

func (node *BurnerNode) tStagIn() float64 {
	return node.GasInput().GetState().(states.GasPortState).TStag
}

func (node *BurnerNode) tStagOut() float64 {
	return node.GasOutput().GetState().(states.GasPortState).TStag
}

func (node *BurnerNode) pStagIn() float64 {
	return node.GasInput().GetState().(states.GasPortState).PStag
}

func (node *BurnerNode) pStagOut() float64 {
	return node.GasOutput().GetState().(states.GasPortState).PStag
}

func (node *BurnerNode) getFuelParameters(initAlpha float64) (float64, float64) {
	var converged = func(currAlpha, nextAlpha float64) bool {
		return math.Abs(currAlpha-nextAlpha)/currAlpha <= node.precision
	}

	var currAlpha = initAlpha
	var nextAlpha = node.getNextAlpha(currAlpha)

	for !converged(currAlpha, nextAlpha) {
		currAlpha = nextAlpha
		nextAlpha = node.getNextAlpha(currAlpha)
	}

	var fuelMassRateRel = node.getFuelMassRateRel(nextAlpha)
	return fuelMassRateRel, nextAlpha
}

func (node *BurnerNode) getNextAlpha(currAlpha float64) float64 {
	return 1 / (node.getFuelMassRateRel(currAlpha) * node.fuel.AirMassTheory())
}

func (node *BurnerNode) getFuelMassRateRel(currAlpha float64) float64 {
	node.outletGas = node.fuel.GetCombustionGas(currAlpha)

	var num1 = gases.CpMean(node.outletGas, node.tgStag, node.t0, defaultN) * (node.tgStag - node.t0)
	var num2 = -gases.CpMean(node.inletGas(), node.tStagIn(), node.t0, defaultN) * (node.tStagIn() - node.t0)

	var denom1 = node.fuel.QLower() * node.etaBurn
	var denom2 = -gases.CpMean(node.outletGas, node.tgStag, node.t0, defaultN) * (node.tgStag - node.t0)
	var denom3 = fuel.CpMean(node.fuel, node.tFuel, node.t0, defaultN) * (node.tFuel - node.t0)

	return (num1 + num2) / (denom1 + denom2 + denom3)
}