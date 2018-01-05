package constructive

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type TurbineCharacteristic func(lambdaU, normPiStag float64) float64

type ParametricBlockedTurbineNode interface {
	TurbineNode
	nodes.PowerSink
	nodes.RPMSink
	NormPiT() float64
	SetNormPiT(normPiT float64)
}

func NewParametricBlockedTurbineNode(
	massRate0, piT0, eta0, t0, p0, inletMeanDiameter, precision float64,
	leakMassRateFunc, coolMasRateRel, inflowMassRateRel func(TurbineNode) float64,
	normMassRateChar, normEtaChar TurbineCharacteristic,
) ParametricBlockedTurbineNode {
	var result = &parametricBlockedTurbineNode{
		precision: precision,

		t0: t0,
		p0: p0,

		massRate0: massRate0,
		piT0:      piT0,
		eta0:      eta0,

		inletMeanDiameter: inletMeanDiameter,

		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateRel,
		inflowMassRateRel: inflowMassRateRel,

		normMassRateChar: normMassRateChar,
		normEtaChar:      normEtaChar,

		normPiT: 1,
	}

	result.baseBlockedTurbine = NewBaseBlockedTurbine(result, precision)
	result.rpmInput = graph.NewAttachedPort(result)
	return result
}

type parametricBlockedTurbineNode struct {
	graph.BaseNode
	*baseBlockedTurbine

	rpmInput graph.Port

	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64

	normMassRateChar TurbineCharacteristic
	normEtaChar      TurbineCharacteristic

	t0 float64
	p0 float64

	massRate0 float64
	piT0      float64
	eta0      float64

	inletMeanDiameter float64
	normPiT           float64

	precision float64
}

func (node *parametricBlockedTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ParametricBlockedTurbine")
}

func (node *parametricBlockedTurbineNode) GetPorts() []graph.Port {
	return append(node.baseBlockedTurbine.GetPorts(), node.rpmInput)
}

// parametric turbine does not declare massRateInput as its required
// port, cos massRate is its inner property which is balanced
// with solver while solving the whole system
func (node *parametricBlockedTurbineNode) GetRequirePorts() []graph.Port {
	return []graph.Port{
		node.baseBlockedTurbine.gasInput,
		node.baseBlockedTurbine.temperatureInput,
		node.baseBlockedTurbine.pressureInput,
		node.baseBlockedTurbine.powerInput,
		node.rpmInput,
	}
}

func (node *parametricBlockedTurbineNode) Process() error {
	var tStagOut, err = node.getTStagOut(node.turbineLabour())
	if err != nil {
		return err
	}

	var piTStag = node.piTStag()
	var pStagOut = node.pStagIn() / piTStag

	var massRateOut = node.massRateInput.GetState().(states.MassRatePortState).MassRate * node.massRateRelFactor()

	graph.SetAll(
		[]graph.PortState{
			node.gasInput.GetState(),
			states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut),
			states.NewMassRatePortState(massRateOut),
			states.NewPowerPortState(node.turbineLabour()), // TODO maybe need to pass sum of labours
		},
		[]graph.Port{
			node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
			node.powerOutput,
		},
	)

	return nil
}

func (node *parametricBlockedTurbineNode) Eta() float64 {
	return node.etaT()
}

func (node *parametricBlockedTurbineNode) NormPiT() float64 {
	return node.normPiT
}

func (node *parametricBlockedTurbineNode) SetNormPiT(normPiT float64) {
	node.normPiT = normPiT
}

func (node *parametricBlockedTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *parametricBlockedTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
}

func (node *parametricBlockedTurbineNode) PiTStag() float64 {
	return node.piTStag()
}

func (node *parametricBlockedTurbineNode) RPMInput() graph.Port {
	return node.rpmInput
}

func (node *parametricBlockedTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}

func (node *parametricBlockedTurbineNode) etaT() float64 {
	return node.normEtaChar(node.lambdaU(), node.normPiT) * node.piT0
}

func (node *parametricBlockedTurbineNode) massRate() float64 {
	return MassRate(
		node.normMassRate(), node.massRate0,
		node.tStagIn(), node.t0, node.pStagIn(), node.p0,
	)
}

func (node *parametricBlockedTurbineNode) normMassRate() float64 {
	return node.normMassRateChar(node.lambdaU(), node.normPiT)
}

func (node *parametricBlockedTurbineNode) lambdaU() float64 {
	var u = math.Pi * node.inletMeanDiameter * node.rpm() / 60

	var r = node.inputGas().R()
	var k = gases.K(node.inputGas(), node.tStagIn())
	var aCrit = gdf.ACrit(k, r, node.tStagIn())
	return u / aCrit
}

func (node *parametricBlockedTurbineNode) rpm() float64 {
	return node.rpmInput.GetState().(states.RPMPortState).RPM
}

func (node *parametricBlockedTurbineNode) piTStag() float64 {
	return node.normPiT * node.piT0
}
