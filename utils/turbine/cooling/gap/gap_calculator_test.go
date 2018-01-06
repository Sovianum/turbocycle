package gap

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GapCalculatorTestSuite struct {
	suite.Suite

	cooler gases.Gas
	gas    gases.Gas

	ca   float64
	pGas float64

	ba          float64
	bladeLength float64
	perimeter   float64
	bladeArea   float64
	wallThk     float64

	lambdaM float64

	nuGasFunc func(re float64) float64

	tGas         float64
	tWallOuter   float64
	tCoolerInlet float64

	gapCalculator gapCalculator
}

func (suite *GapCalculatorTestSuite) SetupTest() {
	suite.cooler = gases.GetAir()
	suite.gas = fuel.GetCH4().GetCombustionGas(3)

	suite.ba = 0.128
	suite.bladeLength = 0.0676
	suite.perimeter = 2 * suite.ba
	suite.bladeArea = suite.bladeLength * suite.perimeter

	suite.wallThk = 1.8e-3

	suite.pGas = 0.5e6

	suite.tGas = 1223
	suite.tWallOuter = 1023
	suite.tCoolerInlet = 600

	var dMean = 0.6467
	var gasMassRate float64 = 439
	var gasDensity = gases.Density(suite.gas, suite.tGas, suite.pGas)
	var alpha1 = common.ToRadians(21 + 1/6)
	var area = math.Pi * dMean * suite.bladeLength
	suite.ca = gasMassRate / (area * gasDensity * math.Sin(alpha1))

	suite.lambdaM = 20

	suite.nuGasFunc = func(re float64) float64 {
		return 0.079 * math.Pow(re, 0.68)
	}

	suite.gapCalculator = gapCalculator{
		cooler:       suite.cooler,
		gas:          suite.gas,
		ca:           suite.ca,
		pGas:         suite.pGas,
		ba:           suite.ba,
		bladeLength:  suite.bladeLength,
		perimeter:    suite.perimeter,
		bladeArea:    suite.bladeArea,
		wallThk:      suite.wallThk,
		lambdaM:      suite.lambdaM,
		nuGasFunc:    suite.nuGasFunc,
		tGas:         suite.tGas,
		tWallOuter:   suite.tWallOuter,
		tCoolerInlet: suite.tCoolerInlet,
	}
}

func (suite *GapCalculatorTestSuite) TestSmoke() {
	var pack = suite.gapCalculator.GetPack(0.05)
	assert.Nil(suite.T(), pack.Err)
}

func TestBladeProfileTestSuite(t *testing.T) {
	suite.Run(t, new(GapCalculatorTestSuite))
}
