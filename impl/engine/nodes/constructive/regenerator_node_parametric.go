package constructive

import (
	"fmt"
	math2 "math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"gonum.org/v1/gonum/mat"
)

const (
	temperaturePrecision = 1e-3
)

type NuFunc func(gas gases.Gas, velocity, pressure, temperature, d float64) float64

func DefaultNuFunc(gas gases.Gas, velocity, pressure, temperature, d float64) float64 {
	var density = gases.Density(gas, temperature, pressure)
	var viscosity = gas.Mu(temperature)
	var re = velocity * d * density / viscosity

	var lambda = gas.Lambda(temperature)
	var cp = gas.Cp(temperature)
	var pr = viscosity * cp / lambda

	return 0.56 * math2.Pow(re, 0.5) * math2.Pow(pr, 0.36)
}

type TemperatureDropFunc func(tHotIn, tHotOut, tColdIn, tColdOut float64) float64

func CounterTDrop(tHotIn, tHotOut, tColdIn, tColdOut float64) float64 {
	var dtHot = tHotIn - tColdOut
	var dtCold = tHotOut - tColdIn
	return logTDrop(dtHot, dtCold)
}

func FrowardTDrop(tHotIn, tHotOut, tColdIn, tColdOut float64) float64 {
	var dtHot = tHotIn - tColdIn
	var dtCold = tHotOut - tColdOut
	if dtCold < 0 {
		dtCold = 1
	}
	return logTDrop(dtHot, dtCold)
}

func logTDrop(dt1, dt2 float64) float64 {
	if math2.Abs(dt1-dt2) < 1e-7 {
		return (dt1 - dt2) / (dt1 / dt2)
	}
	return (dt1 - dt2) / math2.Log(dt1/dt2)
}

func NewParametricRegeneratorNodeFromProto(
	proto RegeneratorNode,
	massRateHot0, massRateCold0,
	velocityHot0, velocityCold0,
	hydraulicDiameterHot, hydraulicDiameterCold,
	precision, relaxCoef float64, iterLimit int,
	meanTemperatureDropFunc TemperatureDropFunc,
	nuHotFunc, nuColdFunc NuFunc,
) RegeneratorNode {
	result := NewParametricRegeneratorNode(
		proto.HotInput().GasInput().GetState().Value().(gases.Gas),
		proto.ColdInput().GasInput().GetState().Value().(gases.Gas),
		massRateHot0, massRateCold0,
		proto.HotInput().TemperatureInput().GetState().Value().(float64),
		proto.ColdInput().TemperatureInput().GetState().Value().(float64),
		proto.HotInput().PressureInput().GetState().Value().(float64),
		proto.ColdInput().PressureInput().GetState().Value().(float64),
		velocityHot0, velocityCold0,
		proto.Sigma(), hydraulicDiameterHot, hydraulicDiameterCold,
		precision, relaxCoef, iterLimit,
		meanTemperatureDropFunc, nuHotFunc, nuColdFunc,
	)
	graph.SetAll(
		[]graph.PortState{
			proto.HotInput().GasInput().GetState(),
			proto.HotInput().TemperatureInput().GetState(),
			proto.HotInput().PressureInput().GetState(),
			proto.HotInput().MassRateInput().GetState(),
			proto.ColdInput().GasInput().GetState(),
			proto.ColdInput().TemperatureInput().GetState(),
			proto.ColdInput().PressureInput().GetState(),
			proto.ColdInput().MassRateInput().GetState(),
		},
		[]graph.Port{
			result.HotInput().GasInput(),
			result.HotInput().TemperatureInput(),
			result.HotInput().PressureInput(),
			result.HotInput().MassRateInput(),
			result.ColdInput().GasInput(),
			result.ColdInput().TemperatureInput(),
			result.ColdInput().PressureInput(),
			result.ColdInput().MassRateInput(),
		},
	)
	return result
}

func NewParametricRegeneratorNode(
	hotGas0, coldGas0 gases.Gas,
	massRateHot0, massRateCold0, tHotIn0, tColdIn0,
	pHotIn0, pColdIn0, velocityHot0, velocityCold0,
	sigma0, hydraulicDiameterHot, hydraulicDiameterCold,
	precision, relaxCoef float64, iterLimit int,
	meanTemperatureDropFunc TemperatureDropFunc,
	nuHotFunc, nuColdFunc NuFunc,
) RegeneratorNode {
	var result = &parametricRegeneratorNode{
		hotGas0:               hotGas0,
		coldGas0:              coldGas0,
		massRateHot0:          massRateHot0,
		massRateCold0:         massRateCold0,
		tHotIn0:               tHotIn0,
		tColdIn0:              tColdIn0,
		pHotIn0:               pHotIn0,
		pColdIn0:              pColdIn0,
		velocityHot0:          velocityHot0,
		velocityCold0:         velocityCold0,
		sigma0:                sigma0,
		hydraulicDiameterHot:  hydraulicDiameterHot,
		hydraulicDiameterCold: hydraulicDiameterCold,
		precision:             precision,
		relaxCoef:             relaxCoef,
		iterLimit:             iterLimit,

		meanTemperatureDropFunc: meanTemperatureDropFunc,
		nuHotFunc:               nuHotFunc,
		nuColdFunc:              nuColdFunc,
	}
	result.baseRegenerator = newBaseRegenerator(result)

	return result
}

// wall heat conductivity is not taken into account
type parametricRegeneratorNode struct {
	*baseRegenerator

	massRateCold0 float64
	massRateHot0  float64
	pHotIn0       float64
	pColdIn0      float64
	tHotIn0       float64
	tColdIn0      float64
	sigma0        float64

	velocityHot0  float64
	velocityCold0 float64

	hotGas0  gases.Gas
	coldGas0 gases.Gas

	hydraulicDiameterHot  float64
	hydraulicDiameterCold float64

	precision float64
	relaxCoef float64
	iterLimit int

	meanTemperatureDropFunc TemperatureDropFunc

	nuHotFunc  NuFunc
	nuColdFunc NuFunc

	heatExchangeArea float64
	hotArea          float64
	coldArea         float64
	sigma            float64
}

func (node *parametricRegeneratorNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "ParametricRegenerator")
}

func (node *parametricRegeneratorNode) Sigma() float64 {
	return node.sigma
}

func (node *parametricRegeneratorNode) Process() error {
	var geomErr = node.setGeomParameters()
	if geomErr != nil {
		return geomErr
	}

	var tHotOut, tColdOut, tErr = node.getOutputTemperatures()
	if tErr != nil {
		return tErr
	}

	var tHotIn, tColdIn = node.tStagHotIn(), node.tStagColdIn()
	node.sigma = (tColdOut - tColdIn) / (tHotIn - tColdIn)

	graph.SetAll(
		[]graph.PortState{
			states.NewTemperaturePortState(tColdOut), states.NewTemperaturePortState(tHotOut),
			node.coldMassRateInput.GetState(), node.hotMassRateInput.GetState(),
			node.coldGasInput.GetState(), node.hotGasInput.GetState(),
			node.coldPressureInput.GetState(), node.hotPressureInput.GetState(),
		},
		[]graph.Port{
			node.coldTemperatureOutput, node.hotTemperatureOutput,
			node.coldMassRateOutput, node.hotMassRateOutput,
			node.coldGasOutput, node.hotGasOutput,
			node.coldPressureOutput, node.hotPressureOutput,
		},
	)
	return nil
}

func (node *parametricRegeneratorNode) getOutputTemperatures() (float64, float64, error) {
	hotMassRate := node.hotMassRateInput.GetState().(states.MassRatePortState).MassRate
	coldMassRate := node.coldMassRateInput.GetState().(states.MassRatePortState).MassRate

	hotGas := node.hotGasInput.GetState().(states.GasPortState).Gas
	coldGas := node.coldGasInput.GetState().(states.GasPortState).Gas

	tHotIn, tColdIn := node.tStagHotIn(), node.tStagColdIn()

	pHot := node.hotPressureInput.GetState().(states.PressurePortState).PStag
	pCold := node.coldPressureInput.GetState().(states.PressurePortState).PStag

	densityHot := gases.Density(hotGas, tHotIn, pHot)
	densityCold := gases.Density(coldGas, tColdIn, pCold)

	cHot := hotMassRate / (node.hotArea * densityHot)
	cCold := coldMassRate / (node.coldArea * densityCold)

	residualFunc := func(tVec *mat.VecDense) (*mat.VecDense, error) {
		tHotOut, tColdOut := tVec.At(0, 0), tVec.At(1, 0)

		cpHot := gases.CpMean(hotGas, tHotIn, tHotOut, nodes.DefaultN)
		qHot := hotMassRate * cpHot * (tHotIn - tHotOut)

		cpCold := gases.CpMean(coldGas, tColdIn, tColdOut, nodes.DefaultN)
		qCold := coldMassRate * cpCold * (tColdOut - tColdIn)

		tDrop := node.meanTemperatureDropFunc(tHotIn, tHotOut, tColdIn, tColdOut)

		tColdMean := tColdIn + tDrop/2
		tHotMean := tHotIn - tDrop/2

		heatTransferCoef := node.getHeatTransferCoef(
			hotGas, coldGas, tHotMean, tColdMean, pHot, pCold, cHot, cCold,
		)

		qTransfer := node.heatExchangeArea * heatTransferCoef * tDrop

		res := mat.NewVecDense(2, []float64{qHot - qCold, qCold - qTransfer})
		if math2.IsNaN(mat.Norm(res, 2)) {
			return nil, fmt.Errorf("NaN obtained")
		}

		return res, nil
	}

	eqSystem := math.NewEquationSystem(residualFunc, 2)
	solver, solverErr := newton.NewUniformNewtonSolver(eqSystem, 1e-3, newton.NoLog)
	if solverErr != nil {
		return 0, 0, solverErr
	}

	dt := tHotIn - tColdIn
	solution, solutionErr := solver.Solve(
		mat.NewVecDense(2, []float64{tHotIn - 0.1*dt, tColdIn + 0.1*dt}), temperaturePrecision, 0.5, 1000,
	)
	if solutionErr != nil {
		return 0, 0, solutionErr
	}

	return solution.At(0, 0), solution.At(1, 0), nil
}

func (node *parametricRegeneratorNode) setGeomParameters() error {
	var meanTDrop0, err = node.getMeanTDrop0()
	if err != nil {
		return err
	}

	var tColdIn0, tColdOut0 = node.tColdIn0, node.getTColdOut0()

	var q0 = node.massRateCold0 * gases.CpMean(node.coldGas0, tColdIn0, tColdOut0, nodes.DefaultN) * (tColdOut0 - tColdIn0)
	var heatExchangeCoef = node.getHeatTransferCoef0(meanTDrop0)

	var heatExchangeArea = q0 / (heatExchangeCoef * meanTDrop0)

	var hotDensity = gases.Density(node.hotGas0, node.tHotIn0, node.pHotIn0)
	var coldDensity = gases.Density(node.coldGas0, node.tColdIn0, node.pColdIn0)

	var hotArea = node.massRateHot0 / (hotDensity * node.velocityHot0)
	var coldArea = node.massRateCold0 / (coldDensity * node.velocityCold0)

	node.heatExchangeArea, node.hotArea, node.coldArea = heatExchangeArea, hotArea, coldArea

	return nil
}

func (node *parametricRegeneratorNode) getHeatTransferCoef0(meanTDrop0 float64) float64 {
	var tColdMean = node.tColdIn0 + meanTDrop0/2
	var tHotMean = node.tHotIn0 - meanTDrop0/2

	return node.getHeatTransferCoef(
		node.hotGas0, node.coldGas0,
		tHotMean, tColdMean, node.pHotIn0, node.pColdIn0,
		node.velocityHot0, node.velocityCold0,
	)
}

func (node *parametricRegeneratorNode) getHeatTransferCoef(
	hotGas, coldGas gases.Gas,
	tHotMean, tColdMean,
	pHot, pCold,
	cHot, cCold float64,
) float64 {
	var nuCold = node.nuColdFunc(coldGas, cCold, pCold, tColdMean, node.hydraulicDiameterCold)
	var nuHot = node.nuHotFunc(hotGas, cHot, pHot, tHotMean, node.hydraulicDiameterHot)

	var lambdaCold = node.coldGas0.Lambda(tColdMean)
	var lambdaHot = node.hotGas0.Lambda(tHotMean)

	var alphaCold = lambdaCold * nuCold / node.hydraulicDiameterCold
	var alphaHot = lambdaHot * nuHot / node.hydraulicDiameterHot

	var k0 = 1 / (1/alphaCold + 1/alphaHot)
	return k0
}

func (node *parametricRegeneratorNode) getMeanTDrop0() (float64, error) {
	var tColdOut0 = node.getTColdOut0()
	var tHotOut0, err = node.getTHotOut0()
	if err != nil {
		return 0, err
	}

	var meanTDrop = node.meanTemperatureDropFunc(node.tHotIn0, tHotOut0, node.tColdIn0, tColdOut0)
	return meanTDrop, nil
}

func (node *parametricRegeneratorNode) getTHotOut0() (float64, error) {
	var cpCold = gases.CpMean(node.coldGas0, node.tColdIn0, node.getTColdOut0(), nodes.DefaultN)

	var iterFunc = func(tHotOut0 float64) (float64, error) {
		var massRateCoef = node.massRateCold0 / node.massRateHot0

		var cpHot = gases.CpMean(node.hotGas0, node.tHotIn0, tHotOut0, nodes.DefaultN)
		var cpCoef = cpCold / cpHot

		return node.tHotIn0 - massRateCoef*cpCoef*node.sigma0*(node.tHotIn0-node.tColdIn0), nil
	}

	return common.SolveIteratively(iterFunc, node.tHotIn0, node.precision, node.relaxCoef, node.iterLimit)
}

func (node *parametricRegeneratorNode) getTColdOut0() float64 {
	return node.tColdIn0 + node.sigma0*(node.tHotIn0-node.tColdIn0)
}
