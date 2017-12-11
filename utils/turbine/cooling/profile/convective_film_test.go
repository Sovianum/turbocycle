package profile

import (
	"testing"

	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/mat"
)

type ConvectiveFilmTestSuite struct {
	suite.Suite

	solver ode.Solver

	coolerMassRate0 float64

	cooler gases.Gas
	gas    gases.Gas

	gasTempStag        func(x float64) float64
	gasPressureStag    func(x float64) float64
	coolerPressureStag func(x float64) float64

	lambdaFunc func(xRel float64) float64

	alphaCoolerConv func(x, theta float64) float64
	alphaGasConv    func(x, theta float64) float64

	slitInfoArray []SlitInfo

	wallThk      func(x float64) float64
	lambdaM      func(t float64) float64
	segment      geom.Segment
	solutionStep float64
	tSystem      TemperatureSystem
}

func (suite *ConvectiveFilmTestSuite) SetupTest() {
	suite.solver = ode.NewEulerSolver()
	suite.coolerMassRate0 = 0.05

	suite.cooler = gases.GetAir()
	suite.gas = fuel.GetCH4().GetCombustionGas(3)

	suite.gasTempStag = func(x float64) float64 {
		return 1400
	}
	suite.gasPressureStag = func(x float64) float64 {
		return 2e6
	}
	suite.coolerPressureStag = func(x float64) float64 {
		return 2e6 * 1.5
	}

	suite.lambdaFunc = func(xRel float64) float64 {
		return 0.9
	}

	suite.alphaCoolerConv = func(x, theta float64) float64 {
		return 1e4
	}
	suite.alphaGasConv = func(x, theta float64) float64 {
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
				mat.NewVecDense(2, []float64{1e-1, 0}),
			},
		),
		0, 1,
	)
	suite.solutionStep = 1e-4

	suite.setTSystem()
}

func (suite *ConvectiveFilmTestSuite) TestFilmAlphaHole() {
	var step = 5e-3
	var t0 float64 = 900

	suite.alphaCoolerConv = func(x, theta float64) float64 {
		return 100
	}
	suite.alphaGasConv = func(x, theta float64) float64 {
		return 150
	}

	var holeCoord = 3e-2
	suite.slitInfoArray = []SlitInfo{
		NewSlitInfo(holeCoord, 1e-4, 25e-8, 1, 1),
		NewSlitInfo(holeCoord*2, 1e-4, 1e-4, 1, 1),
	}

	suite.setTSystem()
	var solution = suite.tSystem.Solve(0, t0, 1, step)

	for i := range solution.ParametricCoord {
		var alphaGas = suite.alphaGasConv(solution.LengthCoord[i], solution.AirTemperature[i])
		var alphaFilm = solution.AlphaGas[i]

		if solution.LengthCoord[i] < holeCoord {
			assert.InDelta(
				suite.T(),
				alphaGas,
				alphaFilm,
				1e-9,
				"i = %d, alpha_gas = %.2f, alpha_film = %.2f",
				i, alphaGas, alphaFilm,
			)
		} else {
			assert.True(
				suite.T(),
				alphaFilm > alphaGas,
				"i = %d, alpha_gas = %.2f, alpha_film = %.2f",
				i, alphaGas, alphaFilm,
			)
		}
	}
}

func TestConvectiveFilmTestSuite(t *testing.T) {
	suite.Run(t, new(ConvectiveFilmTestSuite))
}

func (suite *ConvectiveFilmTestSuite) setTSystem() {
	suite.tSystem = NewConvFilmTemperatureSystem(
		suite.solver,
		suite.coolerMassRate0,
		suite.cooler, suite.gas,
		suite.gasTempStag, suite.gasPressureStag,
		suite.coolerPressureStag,
		suite.lambdaFunc,
		suite.alphaCoolerConv,
		suite.alphaGasConv,
		suite.slitInfoArray,
		suite.wallThk,
		suite.lambdaM,
		suite.segment,
	)
}
