package turbine

import (
	"testing"

	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/ditributions"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

const (
	expectedPi  = 3
	expectedEta = 0.85

	lambdaOut = 0.3
	alphaT    = 3.5
	tBlockedT = 1500
	pBlockedT = 6e5
	lBlockedT = 3.5e5
)

type TurbineFittingTetSuite struct {
	suite.Suite
	node       StagedTurbineNode
	htVariator variator.Variator
	piSolver   variator.VariatorSolver
}

func (suite *TurbineFittingTetSuite) SetupTest() {
	suite.node = getTestTurbine()
	suite.Require().NoError(suite.node.Process())
	suite.htVariator = variator.FromCallables(suite.node.Ht, suite.node.SetHt)
	suite.piSolver = variator.NewVariatorSolver(
		func() (*mat.VecDense, error) {
			if err := suite.node.Process(); err != nil {
				return nil, err
			}
			return mat.NewVecDense(1, []float64{PiStag(suite.node) - expectedPi}), nil
		},
		[]variator.Variator{suite.htVariator},
		newton.NewUniformNewtonSolverGen(1e-5, newton.NoLog),
	)
	suite.Require().NoError(suite.node.Process())
}

func (suite TurbineFittingTetSuite) TestFitPi() {
	_, err := suite.piSolver.Solve(mat.NewVecDense(1, []float64{suite.node.Ht()}), 1e-6, 0.1, 1000)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
}

func (suite TurbineFittingTetSuite) TestFitPiAndEta() {
	eqSys := GetTurbinePiEtaEqSys(
		suite.node,
		common.Scaler(ditributions.GetUnitConstant().Scale(phi)),
		common.Scaler(ditributions.GetUnitConstant().Scale(psi)),
		expectedPi, expectedEta,
	)
	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(2, []float64{suite.node.Ht(), phi}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(expectedPi, PiStag(suite.node), 1e-5)
	suite.InDelta(expectedEta, EtaStag(suite.node), 1e-5)
}

func (suite TurbineFittingTetSuite) TestFitCycle() {
	cycleTurbine := suite.getTestBlockedTurbine()
	eqSys := GetCycleFitEqSys(
		suite.node, cycleTurbine,
		common.Scaler(ditributions.GetUnitConstant().Scale(phi)),
		common.Scaler(ditributions.GetUnitConstant().Scale(psi)),
	)

	solver, err := newton.NewUniformNewtonSolver(eqSys, 1e-5, newton.NoLog)
	suite.Require().NoError(err)

	_, err = solver.Solve(mat.NewVecDense(2, []float64{suite.node.Ht(), phi}), 1e-6, 1, 100)
	suite.Require().NoError(err)

	suite.InDelta(cycleTurbine.PiTStag(), PiStag(suite.node), 1e-5)
	suite.InDelta(expectedEta, EtaStag(suite.node), 1e-5)

	suite.InDelta(suite.node.TemperatureOutput().GetState().Value().(float64), cycleTurbine.TStagOut(), 2)
	suite.InDelta(suite.node.PressureOutput().GetState().Value().(float64), cycleTurbine.PStagOut(), 1)
}

func (suite TurbineFittingTetSuite) getTestBlockedTurbine() constructive.BlockedTurbineNode {
	turbine := constructive.NewBlockedTurbineNode(
		expectedEta, lambdaOut, 0.05,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
	)

	gasState := states.NewGasPortState(fuel.GetCH4().GetCombustionGas(gases.GetAir(), alphaT))
	pressureState := states.NewPressurePortState(pBlockedT)
	temperatureState := states.NewTemperaturePortState(tBlockedT)
	massRateState := states.NewMassRatePortState(1)

	graph.SetAll(
		[]graph.PortState{gasState, pressureState, temperatureState, massRateState},
		[]graph.Port{turbine.GasInput(), turbine.PressureInput(), turbine.TemperatureInput(), turbine.MassRateInput()},
	)

	powerState := states.NewPowerPortState(-lBlockedT)
	turbine.PowerInput().SetState(powerState)

	turbine.Process()
	return turbine
}

func TestCompressorFittingTestSuite(t *testing.T) {
	suite.Run(t, new(TurbineFittingTetSuite))
}
