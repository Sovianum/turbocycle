package geometry

import (
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
)

const (
	lRelOut    = 0.2
	elongation = 4.
	deltaRel   = 0.1
)

type BladingGeometryGeneratorTestSuite struct {
	suite.Suite
	gen BladingGeometryGenerator
}

func (suite *BladingGeometryGeneratorTestSuite) SetupTest() {
	suite.gen = NewGeneratorFromProfileAngles(
		lRelOut, elongation, deltaRel, gammaIn, gammaOut,
	)
}

func (suite *BladingGeometryGeneratorTestSuite) TestAngleConversions() {
	var expectedInnerAngle = -0.5
	var expectedOuterAngle = 0.5
	var expectedMeanAngle = 0.
	var expectedTotalAngle = 1.

	var innerAngle, outerAngle = GetInnerAndOuterAngles(expectedTotalAngle, expectedMeanAngle)
	var totalAngle, meanAngle = GetTotalAndMeanLineAngles(expectedInnerAngle, expectedOuterAngle)

	assert.True(
		suite.T(),
		common.ApproxEqual(expectedInnerAngle, innerAngle, 0.00001),
		testMessage(expectedInnerAngle, innerAngle),
	)
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedOuterAngle, outerAngle, 0.00001),
		testMessage(expectedOuterAngle, outerAngle),
	)
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedTotalAngle, totalAngle, 0.00001),
		testMessage(expectedTotalAngle, totalAngle),
	)
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedMeanAngle, meanAngle, 0.00001),
		testMessage(expectedMeanAngle, meanAngle),
	)
}

func (suite *BladingGeometryGeneratorTestSuite) TestParameters() {
	assert.Equal(suite.T(), lRelOut, suite.gen.LRelOut())
	assert.Equal(suite.T(), elongation, suite.gen.Elongation())
	assert.Equal(suite.T(), gammaIn, suite.gen.GammaIn())
	assert.Equal(suite.T(), gammaOut, suite.gen.GammaOut())
}

func (suite *BladingGeometryGeneratorTestSuite) TestGenerateFromInlet() {
	var geom = suite.gen.GenerateFromInlet(dMeanIn)

	var expectedDMeanIn = (geom.InnerProfile().Diameter(geom.XBladeIn()) + geom.OuterProfile().Diameter(geom.XBladeIn())) / 2
	var expectedDMeanOut = (geom.InnerProfile().Diameter(geom.XGapOut()) + geom.OuterProfile().Diameter(geom.XGapOut())) / 2
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedDMeanIn, geom.MeanProfile().Diameter(0), 0.00001),
	)
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedDMeanOut, geom.MeanProfile().Diameter(geom.XGapOut()), 0.00001),
	)

	var expectedLOut = lRelOut * expectedDMeanOut
	var lOut = Height(geom.XGapOut(), geom)
	assert.True(
		suite.T(),
		common.ApproxEqual(expectedLOut, lOut, 0.0000001),
		testMessage(expectedLOut, lOut),
	)
}

func (suite *BladingGeometryGeneratorTestSuite) TestDiameterContinuity() {
	var statorGeom = suite.gen.GenerateFromInlet(dMeanIn)

	var lRotorIn = Height(statorGeom.XGapOut(), statorGeom)
	var lRotorOut = 1 / (1 - (math.Tan(gammaOut)-math.Tan(gammaIn))/elongation*(1+deltaRel)) * lRotorIn
	var bladeWidth = lRotorOut / elongation
	var bladingWidth = bladeWidth * (1 + deltaRel)
	var dMeanOut = statorGeom.MeanProfile().Diameter(statorGeom.XGapOut() + bladingWidth)
	var lRelOut = lRotorOut / dMeanOut

	assert.True(
		suite.T(),
		common.ApproxEqual(Height(statorGeom.XGapOut(), statorGeom), lRotorIn, 0.000001),
		testMessage(Height(statorGeom.XGapOut(), statorGeom), lRotorIn),
	)

	assert.True(
		suite.T(),
		common.ApproxEqual(Height(statorGeom.XGapOut()+bladingWidth, statorGeom), lRotorOut, 0.000001),
		testMessage(Height(statorGeom.XGapOut()+bladingWidth, statorGeom), lRotorOut),
	)

	assert.True(
		suite.T(),
		common.ApproxEqual(RelativeHeight(statorGeom.XGapOut()+bladingWidth, statorGeom), lRelOut, 0.000001),
		testMessage(RelativeHeight(statorGeom.XGapOut()+bladingWidth, statorGeom), lRelOut),
	)

	var rotorGeomGen = NewGeneratorFromProfileAngles(
		lRelOut, elongation, deltaRel, gammaIn, gammaOut,
	)
	var rotorGeom = rotorGeomGen.GenerateFromInlet(statorGeom.MeanProfile().Diameter(statorGeom.XGapOut()))

	assert.True(
		suite.T(),
		common.ApproxEqual(Height(statorGeom.XGapOut(), statorGeom), Height(0, rotorGeom), 0.00001),
		testMessage(Height(statorGeom.XGapOut(), statorGeom), Height(0, rotorGeom)),
	)

	assert.True(
		suite.T(),
		common.ApproxEqual(bladeWidth, rotorGeom.XBladeOut(), 0.000001),
		testMessage(bladeWidth, rotorGeom.XBladeOut()),
	)

	var dMeanStatorOut = statorGeom.MeanProfile().Diameter(statorGeom.XGapOut())
	var dMeanRotorIn = rotorGeom.MeanProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dMeanStatorOut, dMeanRotorIn, 0.000001),
		testMessage(dMeanStatorOut, dMeanRotorIn),
	)

	var heightRelStatorOut = RelativeHeight(statorGeom.XGapOut(), statorGeom)
	var heightRelRotorIn = RelativeHeight(0, rotorGeom)
	assert.True(
		suite.T(),
		common.ApproxEqual(heightRelStatorOut, heightRelRotorIn, 0.000001),
		testMessage(heightRelStatorOut, heightRelRotorIn),
	)

	var dInnerStatorOut = statorGeom.InnerProfile().Diameter(statorGeom.XGapOut())
	var dInnerRotorIn = rotorGeom.InnerProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dInnerStatorOut, dInnerRotorIn, 0.000001),
		testMessage(dInnerStatorOut, dInnerRotorIn),
	)

	var dOuterStatorOut = statorGeom.OuterProfile().Diameter(statorGeom.XGapOut())
	var dOuterRotorIn = rotorGeom.OuterProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dOuterStatorOut, dOuterRotorIn, 0.000001),
		testMessage(dOuterStatorOut, dOuterRotorIn),
	)
}

func (suite *BladingGeometryGeneratorTestSuite) TestGenerateFromOutlet() {
	var dMeanOut = dMeanIn
	var geom = suite.gen.GenerateFromOutlet(dMeanOut)
	assert.Equal(suite.T(), dMeanOut, geom.MeanProfile().Diameter(geom.XGapOut()))

	var expectedLOut = lRelOut * dMeanOut
	var lOut = Height(geom.XGapOut(), geom)
	assert.InDelta(suite.T(), expectedLOut, lOut, 0.0000001, testMessage(expectedLOut, lOut))
}

func TestBladingGeometryGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(BladingGeometryGeneratorTestSuite))
}
