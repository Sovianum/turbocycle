package constructive

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type TurbineCharFunc func(lambdaU, normPiStag float64) float64

type ParametricTurbineNode interface {
	TurbineNode
	nodes.RPMSink
	nodes.MassRateSink
	NormPiT() float64
	SetNormPiT(normPiT float64)
}

func NewSimpleParametricTurbineNode(
	massRate0, piT0, eta0, t0, p0, inletMeanDiameter, precision,
	leakMassRateCoef, coolMasRateCoef, inflowMassRateCoef float64,
	normMassRateChar, normEtaChar TurbineCharFunc,
) ParametricTurbineNode {
	return NewParametricTurbineNode(
		massRate0, piT0, eta0, t0, p0, inletMeanDiameter, precision,
		func(TurbineNode) float64 {
			return leakMassRateCoef
		},
		func(TurbineNode) float64 {
			return coolMasRateCoef
		},
		func(TurbineNode) float64 {
			return inflowMassRateCoef
		},
		normMassRateChar, normEtaChar,
	)
}

func NewParametricTurbineNode(
	massRate0, piT0, eta0, t0, p0, inletMeanDiameter, precision float64,
	leakMassRateFunc, coolMasRateFunc, inflowMassRateFunc func(TurbineNode) float64,
	normMassRateChar, normEtaChar TurbineCharFunc,
) ParametricTurbineNode {
	var result = &parametricTurbineNode{
		precision: precision,

		t0: t0,
		p0: p0,

		massRate0: massRate0,
		piT0:      piT0,
		eta0:      eta0,

		inletMeanDiameter: inletMeanDiameter,

		leakMassRateFunc:  leakMassRateFunc,
		coolMasRateRel:    coolMasRateFunc,
		inflowMassRateRel: inflowMassRateFunc,

		normMassRateChar: normMassRateChar,
		normEtaChar:      normEtaChar,

		normPiT: 1,
	}

	result.baseBlockedTurbine = NewBaseBlockedTurbine(result, precision)
	result.rpmInput = graph.NewAttachedPortWithTag(result, nodes.RPMInputTag)
	result.massRateInput = graph.NewAttachedPortWithTag(result, nodes.MassRateInputTag)
	return result
}

type parametricTurbineNode struct {
	graph.BaseNode
	*baseBlockedTurbine

	rpmInput      graph.Port
	massRateInput graph.Port

	leakMassRateFunc  func(TurbineNode) float64
	coolMasRateRel    func(TurbineNode) float64
	inflowMassRateRel func(TurbineNode) float64

	normMassRateChar TurbineCharFunc
	normEtaChar      TurbineCharFunc

	t0 float64
	p0 float64

	massRate0 float64
	piT0      float64
	eta0      float64

	inletMeanDiameter float64
	normPiT           float64

	precision float64
}

func (node *parametricTurbineNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ParametricBlockedTurbine")
}

func (node *parametricTurbineNode) GetPorts() []graph.Port {
	return append(node.baseBlockedTurbine.GetPorts(), node.rpmInput, node.massRateInput)
}

// parametric turbine does not declare massRateInput as its required
// port, cos massRate is its inner property which is balanced
// with solver while solving the whole system
func (node *parametricTurbineNode) GetUpdatePorts() ([]graph.Port, error) {
	var ports, err = node.baseBlockedTurbine.GetUpdatePorts()
	if err != nil {
		return nil, err
	}
	return append(ports, node.massRateInput), nil
}

func (node *parametricTurbineNode) GetRequirePorts() ([]graph.Port, error) {
	var ports, err = node.baseBlockedTurbine.GetRequirePorts()
	if err != nil {
		return nil, err
	}
	return append(ports, node.rpmInput), nil
}

func (node *parametricTurbineNode) Process() error {
	var tStagOut, err = node.getTStagOut()
	if err != nil {
		return err
	}

	var piTStag = node.piTStag()
	var pStagOut = node.pStagIn() / piTStag

	var massRateIn = node.massRate()
	var massRateOut = massRateIn * node.massRateRelFactor()

	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), tStagOut, nodes.DefaultN)
	var lSpecific = cp * (node.tStagIn() - tStagOut)

	graph.SetAll(
		[]graph.PortState{
			states.NewMassRatePortState(massRateIn),

			node.gasInput.GetState(),
			states.NewTemperaturePortState(tStagOut),
			states.NewPressurePortState(pStagOut),
			states.NewMassRatePortState(massRateOut),
			states.NewPowerPortState(lSpecific),
		},
		[]graph.Port{
			node.massRateInput,
			node.gasOutput, node.temperatureOutput, node.pressureOutput, node.massRateOutput,
			node.powerOutput,
		},
	)

	return nil
}

func (node *parametricTurbineNode) Eta() float64 {
	return node.etaT()
}

func (node *parametricTurbineNode) NormPiT() float64 {
	return node.normPiT
}

func (node *parametricTurbineNode) SetNormPiT(normPiT float64) {
	node.normPiT = normPiT
}

func (node *parametricTurbineNode) LeakMassRateRel() float64 {
	return node.leakMassRateFunc(node)
}

func (node *parametricTurbineNode) CoolMassRateRel() float64 {
	return node.coolMasRateRel(node)
}

func (node *parametricTurbineNode) LSpecific() float64 {
	var cp = gases.CpMean(node.inputGas(), node.tStagIn(), node.tStagOut(), nodes.DefaultN)
	return -cp * (node.tStagIn() - node.tStagOut())
}

func (node *parametricTurbineNode) PiTStag() float64 {
	return node.piTStag()
}

func (node *parametricTurbineNode) RPMInput() graph.Port {
	return node.rpmInput
}

func (node *parametricTurbineNode) MassRateInput() graph.Port {
	return node.massRateInput
}

func (node *parametricTurbineNode) getTStagOut() (float64, error) {
	var t0, err = node.getNewTtStag(0.8 * node.tStagIn()) // TODO move 0.8 out of code
	if err != nil {
		return t0, err
	}
	return common.SolveIterativly(node.getNewTtStag, t0, node.precision, nodes.DefaultN)
}

func (node *parametricTurbineNode) getNewTtStag(currTtStag float64) (float64, error) {
	var k = gases.KMean(node.inputGas(), node.tStagIn(), currTtStag, nodes.DefaultN)
	var tStagIn = node.tStagIn()
	var result = tStagIn * (1 - (1-math.Pow(node.piTStag(), (1-k)/k))*node.etaT())
	if math.IsNaN(result) {
		return 0, fmt.Errorf("nan obtained while calculating tStagOut")
	}
	return result, nil
}

func (node *parametricTurbineNode) massRateRelFactor() float64 {
	return 1 + node.leakMassRateFunc(node) + node.coolMasRateRel(node) + node.inflowMassRateRel(node)
}

func (node *parametricTurbineNode) etaT() float64 {
	return node.normEtaChar(node.lambdaU(), node.normPiT) * node.eta0
}

func (node *parametricTurbineNode) massRate() float64 {
	return MassRate(
		node.normMassRate(), node.massRate0,
		node.tStagIn(), node.t0, node.pStagIn(), node.p0,
	)
}

func (node *parametricTurbineNode) normMassRate() float64 {
	return node.normMassRateChar(node.lambdaU(), node.normPiT)
}

func (node *parametricTurbineNode) lambdaU() float64 {
	var u = math.Pi * node.inletMeanDiameter * node.rpm() / 60

	var r = node.inputGas().R()
	var k = gases.K(node.inputGas(), node.tStagIn())
	var aCrit = gdf.ACrit(k, r, node.tStagIn())
	return u / aCrit
}

func (node *parametricTurbineNode) rpm() float64 {
	return node.rpmInput.GetState().(states.RPMPortState).RPM
}

func (node *parametricTurbineNode) piTStag() float64 {
	return node.normPiT * node.piT0
}
