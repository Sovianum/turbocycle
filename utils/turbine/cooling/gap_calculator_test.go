package cooling

import (
	"fmt"
	"math"
	"testing"

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

	suite.ca = 100
	suite.pGas = 0.5e6

	suite.ba = 0.1
	suite.bladeLength = 0.1
	suite.perimeter = 2 * suite.ba
	suite.bladeArea = suite.bladeLength * suite.perimeter
	suite.wallThk = 1e-3

	suite.lambdaM = 20

	suite.nuGasFunc = func(re float64) float64 {
		return 0.079 * math.Pow(re, 0.68)
	}

	suite.tGas = 1300
	suite.tWallOuter = 900
	suite.tCoolerInlet = 500

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
	var pack = suite.gapCalculator.GetPack(1)
	assert.Nil(suite.T(), pack.err)
	fmt.Println(pack)
}

func TestBladeProfileTestSuite(t *testing.T) {
	suite.Run(t, new(GapCalculatorTestSuite))
}
