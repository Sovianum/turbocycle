package geometry

import (
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/Sovianum/turbocycle/common"
	"testing"
	"fmt"
	"math"
)

const (
	bladeWidth = 40e-3
	gapWidth   = 10e-3
	gammaIn    = -0.09
	gammaOut   = 0.06
	dIn        = 1
	dMeanIn    = 1.1
	dOut       = 1.2

)

type BladingGeometryTestSuite struct {
	suite.Suite
	geom BladingGeometry
}

func (suite *BladingGeometryTestSuite) SetupTest() {
	var innerLine = NewAxialProfileLine(0, dIn, gammaIn)
	var outerLine = NewAxialProfileLine(0, dOut, gammaOut)
	suite.geom = NewBladingGeometry(
		bladeWidth, gapWidth,
		innerLine, outerLine,
	)
}

func (suite *BladingGeometryTestSuite) TestMeanLine() {
	var meanAngle = (gammaOut + gammaIn) / 2
	assert.True(
		suite.T(),
		common.ApproxEqual(meanAngle, suite.geom.MeanProfile().Angle(), 0.001),
		testMessage(meanAngle, suite.geom.MeanProfile().Angle()),
	)
	assert.True(suite.T(), common.ApproxEqual(dMeanIn, suite.geom.MeanProfile().Diameter(0), 0.001))

	var xOut = suite.geom.XGapOut()
	var dOut = dMeanIn + 2 * math.Tan(meanAngle) * xOut
	assert.True(
		suite.T(),
		common.ApproxEqual(dOut, suite.geom.MeanProfile().Diameter(xOut), 0.0001),
		testMessage(dOut, suite.geom.MeanProfile().Diameter(xOut)),
	)

}

func (suite *BladingGeometryTestSuite) TestExpansionAngle() {
	assert.True(suite.T(), common.ApproxEqual(gammaOut - gammaIn, ExpansionAngle(suite.geom), 0.001))
}

func TestBladingGeometryTestSuite(t *testing.T) {
	suite.Run(t, new(BladingGeometryTestSuite))
}

func testMessage(x0, x float64) string {
	return fmt.Sprintf("expected: %v; got %v", x0, x)
}

