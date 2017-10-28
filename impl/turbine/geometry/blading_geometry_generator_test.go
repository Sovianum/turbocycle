package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	lRelOut = 0.2
	elongation = 4.
	deltaRel = 0.1
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

func (suite *BladingGeometryGeneratorTestSuite) TestParameters() {
	assert.Equal(suite.T(), lRelOut, suite.gen.LRelOut())
	assert.Equal(suite.T(), elongation, suite.gen.Elongation())
	assert.Equal(suite.T(), gammaIn, suite.gen.GammaIn())
	assert.Equal(suite.T(), gammaOut, suite.gen.GammaOut())
}

func (suite *BladingGeometryGeneratorTestSuite) TestGenerateFromInlet() {
	var geom = suite.gen.GenerateFromInlet(dMeanIn)
	assert.Equal(suite.T(), dMeanIn, geom.MeanProfile().Diameter(0))

	var dMeanOut = geom.MeanProfile().Diameter(geom.XGapOut())
	var expectedLOut = lRelOut * dMeanOut
	var lOut = Height(geom.XGapOut(), geom)
	assert.InDelta(suite.T(), expectedLOut, lOut, 0.0000001, testMessage(expectedLOut, lOut))
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