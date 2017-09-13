package nodes

import (
	"errors"
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/fuel"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
)

type BurnerNode interface {
	core.Node
	GasChannel
	Alpha() float64
	GetFuelRateRel() float64
}

type burnerNode struct {
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
) BurnerNode {
	var result = &burnerNode{
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
	result.ports[gasInput].SetInnerNode(result)

	result.ports[gasOutput] = core.NewPort()
	result.ports[gasOutput].SetInnerNode(result)

	return result
}

func (node *burnerNode) ContextDefined() bool {
	return true
}

func (node *burnerNode) GetPortByTag(tag string) (core.Port, error) {
	switch tag {
	case gasInput:
		return node.gasInput(), nil
	case gasOutput:
		return node.gasOutput(), nil
	default:
		return nil, errors.New(fmt.Sprintf("port with tag \"%s\" not found", tag))
	}
}

func (node *burnerNode) GetRequirePortTags() ([]string, error) {
	return []string{gasInput}, nil
}

func (node *burnerNode) GetUpdatePortTags() ([]string, error) {
	return []string{gasOutput}, nil
}

func (node *burnerNode) GetPortTags() []string {
	return []string{gasInput, gasOutput}
}

func (node *burnerNode) GetPorts() core.PortsType {
	return node.ports
}

func (node *burnerNode) GasInput() core.Port {
	return node.gasInput()
}

func (node *burnerNode) GasOutput() core.Port {
	return node.gasOutput()
}

func (node *burnerNode) Alpha() float64 {
	return node.alpha
}

func (node *burnerNode) TStagIn() float64 {
	return node.tStagIn()
}

func (node *burnerNode) TStagOut() float64 {
	return node.tStagOut()
}

func (node *burnerNode) PStagIn() float64 {
	return node.pStagIn()
}

func (node *burnerNode) PStagOut() float64 {
	return node.pStagOut()
}

func (node *burnerNode) GetFuelRateRel() float64 {
	return node.getFuelMassRateRel(node.alpha)
}

func (node *burnerNode) Process() error {
	var fuelMassRateRel, alpha = node.getFuelParameters(node.initAlpha)
	node.alpha = alpha

	var gasState = node.GasInput().GetState().(states.GasPortState)
	gasState.Gas = node.outletGas
	gasState.TStag = node.tgStag
	gasState.PStag = node.pStagIn() * node.sigma
	gasState.MassRateRel *= 1 + fuelMassRateRel

	node.gasOutput().SetState(gasState)

	return nil
}

func (node *burnerNode) getFuelParameters(initAlpha float64) (float64, float64) {
	var currAlpha = initAlpha
	var nextAlpha = node.getNextAlpha(currAlpha)

	for !common.Converged(currAlpha, nextAlpha, node.precision) {
		currAlpha = nextAlpha
		nextAlpha = node.getNextAlpha(currAlpha)
	}

	var fuelMassRateRel = node.getFuelMassRateRel(nextAlpha)
	return fuelMassRateRel, nextAlpha
}

func (node *burnerNode) getNextAlpha(currAlpha float64) float64 {
	return 1 / (node.getFuelMassRateRel(currAlpha) * node.fuel.AirMassTheory())
}

func (node *burnerNode) getFuelMassRateRel(currAlpha float64) float64 {
	node.outletGas = node.fuel.GetCombustionGas(currAlpha)

	var num1 = gases.CpMean(node.outletGas, node.tgStag, node.t0, defaultN) * (node.tgStag - node.t0)
	var num2 = -gases.CpMean(node.inletGas(), node.tStagIn(), node.t0, defaultN) * (node.tStagIn() - node.t0)

	var denom1 = node.fuel.QLower() * node.etaBurn
	var denom2 = -gases.CpMean(node.outletGas, node.tgStag, node.t0, defaultN) * (node.tgStag - node.t0)
	var denom3 = fuel.CpMean(node.fuel, node.tFuel, node.t0, defaultN) * (node.tFuel - node.t0)

	return (num1 + num2) / (denom1 + denom2 + denom3)
}

func (node *burnerNode) inletGas() gases.Gas {
	return node.gasInput().GetState().(states.GasPortState).Gas
}

func (node *burnerNode) tStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).TStag
}

func (node *burnerNode) tStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).TStag
}

func (node *burnerNode) pStagIn() float64 {
	return node.gasInput().GetState().(states.GasPortState).PStag
}

func (node *burnerNode) pStagOut() float64 {
	return node.gasOutput().GetState().(states.GasPortState).PStag
}

func (node *burnerNode) gasInput() core.Port {
	return node.ports[gasInput]
}

func (node *burnerNode) gasOutput() core.Port {
	return node.ports[gasOutput]
}
