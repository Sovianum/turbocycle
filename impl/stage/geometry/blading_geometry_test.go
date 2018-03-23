package geometry

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	bladeWidth = 40e-3
	gapWidth   = 10e-3
	gammaIn    = -0.09
	gammaOut   = 0.06
	dIn        = 1.
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

func (suite *BladingGeometryTestSuite) TestFreeFunctions() {
	var x = 10.0
	assert.True(
		suite.T(),
		common.ApproxEqual(
			Height(x, suite.geom),
			(suite.geom.OuterProfile().Diameter(x)-suite.geom.InnerProfile().Diameter(x))/2,
			0.000001,
		),
	)

	var lIn = (dOut - dIn) / 2
	var lRelIn = lIn / dMeanIn
	assert.True(suite.T(), common.ApproxEqual(lRelIn, RelativeHeight(0, suite.geom), 0.0001))
	assert.True(suite.T(), common.ApproxEqual(bladeWidth, ChordProjection(suite.geom), 0.00001))

	var elongationIn = lIn / bladeWidth
	assert.True(suite.T(), common.ApproxEqual(elongationIn, Elongation(0, suite.geom), 0.000001))
	assert.True(suite.T(), common.ApproxEqual(gapWidth, AxialGapProjection(suite.geom), 0.000001))

	var expectedArea = math.Pi * lIn * dMeanIn
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedArea, Area(0, suite.geom), 0.00001),
		testMessage(expectedArea, Area(0, suite.geom)),
	)
}

func (suite *BladingGeometryTestSuite) TestMeanLine() {
	var meanAngle = math.Atan((math.Tan(gammaOut) + math.Tan(gammaIn)) / 2)
	assert.True(
		suite.T(),
		common.ApproxEqual(meanAngle, suite.geom.MeanProfile().Angle(), 0.00001),
		testMessage(meanAngle, suite.geom.MeanProfile().Angle()),
	)
	assert.True(suite.T(), common.ApproxEqual(dMeanIn, suite.geom.MeanProfile().Diameter(0), 0.00001))

	var xOut = suite.geom.XGapOut()
	var dOut = dMeanIn + 2*math.Tan(meanAngle)*xOut
	assert.True(
		suite.T(),
		common.ApproxEqual(dOut, suite.geom.MeanProfile().Diameter(xOut), 0.000001),
		testMessage(dOut, suite.geom.MeanProfile().Diameter(xOut)),
	)
}

func (suite *BladingGeometryTestSuite) TestExpansionAngle() {
	assert.True(suite.T(), common.ApproxEqual(gammaOut-gammaIn, ExpansionAngle(suite.geom), 0.00001))
}

func TestBladingGeometryTestSuite(t *testing.T) {
	suite.Run(t, new(BladingGeometryTestSuite))
}

func testMessage(x0, x float64) string {
	return fmt.Sprintf("expected: %v; got %v", x0, x)
}
