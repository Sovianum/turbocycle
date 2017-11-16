package cooling

import (
	"testing"

	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

type TemperatureProfileTestSuite struct {
	suite.Suite

	solver      ode.Solver
	airMassRate float64

	cpAir   func(theta float64) float64
	gasTemp func(x float64) float64

	alphaAir func(x, theta float64) float64
	alphaGas func(x, theta float64) float64

	wallThk      func(x float64) float64
	lambdaM      func(t float64) float64
	segment      geom.Segment
	solutionStep float64
	tSystem      TemperatureSystem
}

func (suite *TemperatureProfileTestSuite) SetupTest() {
	suite.solver = ode.NewEulerSolver()
	suite.airMassRate = 0.05
	suite.cpAir = gases.GetAir().Cp
	suite.gasTemp = func(x float64) float64 {
		return 1400
	}
	suite.alphaAir = func(x, theta float64) float64 {
		return 1e4
	}
	suite.alphaGas = func(x, theta float64) float64 {
		return 1e2
	}

	suite.wallThk = func(x float64) float64 {
		return 1e-3
	}
	suite.lambdaM = func(t float64) float64 {
		return 20
	}
	suite.segment = geom.NewUnitSegment(
		geom.NewBezier(
			[]*mat.VecDense{
				mat.NewVecDense(2, []float64{0, 0}),
				mat.NewVecDense(2, []float64{1, 1}),
			},
		),
		0, 1,
	)
	suite.solutionStep = 1e-4

	suite.setTSystem()
}

func (suite *TemperatureProfileTestSuite) TestHeatTransferQualitative() {
	var step = 1e-1
	var t0 float64 = 900

	suite.alphaAir = func(x, theta float64) float64 {
		return 1
	}

	suite.alphaGas = func(x, theta float64) float64 {
		return 1e9
	}
	suite.setTSystem()
	var solution1 = suite.tSystem.Solve(0, t0, 1, step)

	suite.alphaGas = func(x, theta float64) float64 {
		return 1e8
	}
	suite.setTSystem()
	var solution2 = suite.tSystem.Solve(0, t0, 1, step)

	for i := 0; i < len(solution1.X); i++ {
		assert.InDelta(suite.T(), solution1.HeatTransferCoef[i], solution2.HeatTransferCoef[i], 1e-6)
	}

	suite.alphaGas = func(x, theta float64) float64 {
		return 1
	}

	suite.alphaAir = func(x, theta float64) float64 {
		return 1e8
	}
	suite.setTSystem()
	var solution3 = suite.tSystem.Solve(0, t0, 1, step)

	suite.alphaAir = func(x, theta float64) float64 {
		return 1e9
	}
	suite.setTSystem()
	var solution4 = suite.tSystem.Solve(0, t0, 1, step)

	for i := 0; i < len(solution3.X); i++ {
		assert.InDelta(suite.T(), solution3.HeatTransferCoef[i], solution4.HeatTransferCoef[i], 1e-6)
	}
}

func (suite *TemperatureProfileTestSuite) TestWallTemperatureQualitative() {
	var step = 1e-1
	var t0 float64 = 900

	suite.alphaAir = func(x, theta float64) float64 {
		return 1
	}
	suite.alphaGas = func(x, theta float64) float64 {
		return 1e9
	}
	suite.setTSystem()
	var solution1 = suite.tSystem.Solve(0, t0, 1, step)
	for i := 0; i < len(solution1.X); i++ {
		assert.InDelta(
			suite.T(),
			suite.tSystem.GasTemp(0),
			solution1.WallTemperature[i],
			1e-6,
		)
	}
}

func TestTemperatureProfileTestSuite(t *testing.T) {
	suite.Run(t, new(TemperatureProfileTestSuite))
}

func (suite *TemperatureProfileTestSuite) setTSystem() {
	suite.tSystem = NewTemperatureSystem(
		suite.solver, suite.airMassRate, suite.cpAir,
		suite.gasTemp, suite.alphaAir, suite.alphaGas,
		suite.wallThk, suite.lambdaM, suite.segment,
	)
}