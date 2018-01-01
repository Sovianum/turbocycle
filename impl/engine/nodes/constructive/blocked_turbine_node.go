package constructive

import (
	"errors"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type BlockedTurbineNode interface {
	TurbineNode
	nodes.ComplexGasSource
	nodes.PowerSink
}

func NewBlockedTurbineNode(
	etaT, lambdaOut, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(TurbineNode) float64,
) BlockedTurbineNode {
	var result = &blockedTurbineNode{
		etaT:              etaT,
		precision:         precision,
		lambdaOut:         lambdaOut,
		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateRel,
		inflowMassRateRel: inflowMassRateRel,
	}

	result.powerInput = core.NewAttachedPort(result)
	result.powerOutput = core.NewAttachedPort(result)
	result.complexGasInput = core.NewAttachedPort(result)
	result.complexGasOutput = core.NewAttachedPort(result)

	return result
}

type blockedTurbineNode struct {
	core.BaseNode

	powerInput  core.Port
	powerOutput core.Port

	complexGasInput  core.Port
	complexGasOutput core.Port

	etaT              float64
	precision         float64
	lambdaOut         float64
	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64
}

func (node *blockedTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "BlockedTurbine")
}

func (node *blockedTurbineNode) GetPorts() []core.Port {
	return []core.Port{node.powerInput, node.powerOutput, node.complexGasInput, node.complexGasOutput}
}

func (node *blockedTurbineNode) GetRequirePorts() []core.Port {
	return []core.Port{node.powerInput, node.complexGasInput}
}

func (node *blockedTurbineNode) GetUpdatePorts() []core.Port {
	return []core.Port{node.powerOutput, node.complexGasOutput}
}

func (node *blockedTurbineNode) Process() error {
	var gasState = node.complexGasInput.GetState().(states.ComplexGasPortState)

	var err error
	gasState.TStag, err = node.getTStagOut(node.turbineLabour())
	if err != nil {
		return err
	}

	var piTStag = node.piTStag(gasState.TStag)
	var pi = gdf.Pi(node.lambdaOut, gases.KMean(node.inputGas(), node.tStagIn(), gasState.TStag, nodes.DefaultN))
	gasState.PStag = node.pStagIn() / (piTStag * pi)
	gasState.MassRateRel *= node.massRateRelFactor()

	node.complexGasOutput.SetState(gasState)
	node.powerOutput.SetState(states.NewPowerPortState(node.turbineLabour())) // TODO maybe need to pass sum of labours

	return nil
}

func (node *blockedTurbineNode) LambdaOut() float64 {
	return node.lambdaOut
}

func (node *blockedTurbineNode) Eta() float64 {
	return node.etaT
}

func (node *blockedTurbineNode) LSpecific() float64 {
	return node.turbineLabour()
}

func (node *blockedTurbineNode) PStatOut() float64 {
	return node.pStatOut()
}

func (node *blockedTurbineNode) TStatOut() float64 {
	return node.tStatOut()
}

func (node *blockedTurbineNode) MassRateRel() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).MassRateRel
}

func (node *blockedTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *blockedTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
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
	return node.complexGasInput
}

func (node *blockedTurbineNode) ComplexGasOutput() core.Port {
	return node.complexGasOutput
}

func (node *blockedTurbineNode) PowerInput() core.Port {
	return node.powerInput
}

func (node *blockedTurbineNode) PowerOutput() core.Port {
	return node.powerOutput
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
	return node.complexGasInput.GetState().(states.ComplexGasPortState).Gas
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
	return -node.powerInput.GetState().(states.PowerPortState).LSpecific
}

func (node *blockedTurbineNode) tStatOut() float64 {
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return tStagOut * gdf.Tau(node.lambdaOut, k)
}

func (node *blockedTurbineNode) pStatOut() float64 {
	var pStagOut = node.pStagOut()
	var tStagOut = node.tStagOut()
	var k = gases.K(node.inputGas(), tStagOut)
	return pStagOut * gdf.Tau(node.lambdaOut, k)
}

func (node *blockedTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}

func (node *blockedTurbineNode) tStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).TStag
}

func (node *blockedTurbineNode) pStagIn() float64 {
	return node.complexGasInput.GetState().(states.ComplexGasPortState).PStag
}

func (node *blockedTurbineNode) tStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).TStag
}

func (node *blockedTurbineNode) pStagOut() float64 {
	return node.complexGasOutput.GetState().(states.ComplexGasPortState).PStag
}
