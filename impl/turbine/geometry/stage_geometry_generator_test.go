package geometry

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/Sovianum/turbocycle/common"
)

const (
	rotorLRelOut = 0.2
	rotorElongation = 2
	rotorDeltaRel = 0.1
	rotorGammaIn = -0.09
	rotorGammaOut = 0.09

	statorElongation = 2
	statorDeltaRel = 0.1
	statorGammaIn = -0.09
	statorGammaOut = 0.09

	statorDMeanIn = 1.
)

type StageGeometryGeneratorTestSuite struct {
	suite.Suite
	gen StageGeometryGenerator
}

func (suite *StageGeometryGeneratorTestSuite) SetupTest() {
	var rotorIncompleteGen = NewIncompleteGeneratorFromProfileAngles(
		rotorElongation, rotorDeltaRel, rotorGammaIn, rotorGammaOut,
	)
	var statorIncompleteGen = NewIncompleteGeneratorFromProfileAngles(
		statorElongation, statorDeltaRel, statorGammaIn, statorGammaOut,
	)
	suite.gen = NewStageGeometryGenerator(
		rotorLRelOut, statorIncompleteGen, rotorIncompleteGen,
	)
}

func (suite *StageGeometryGeneratorTestSuite) TestDiameterContinuity() {
	var stageGeom = suite.gen.GenerateFromStatorInlet(statorDMeanIn)

	var dMeanStatorOut = stageGeom.StatorGeometry().MeanProfile().Diameter(stageGeom.StatorGeometry().XGapOut())
	var dMeanRotorIn = stageGeom.RotorGeometry().MeanProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dMeanStatorOut, dMeanRotorIn, 0.000001),
		testMessage(dMeanStatorOut, dMeanRotorIn),
	)

	var heightRelStatorOut = RelativeHeight(stageGeom.StatorGeometry().XGapOut(), stageGeom.StatorGeometry())
	var heightRelRotorIn = RelativeHeight(0, stageGeom.RotorGeometry())
	assert.True(
		suite.T(),
		common.ApproxEqual(heightRelStatorOut, heightRelRotorIn, 0.000001),
		testMessage(heightRelStatorOut, heightRelRotorIn),
	)

	var dInnerStatorOut = stageGeom.StatorGeometry().InnerProfile().Diameter(stageGeom.StatorGeometry().XGapOut())
	var dInnerRotorIn = stageGeom.RotorGeometry().InnerProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dInnerStatorOut, dInnerRotorIn, 0.000001),
		testMessage(dInnerStatorOut, dInnerRotorIn),
	)

	var dOuterStatorOut = stageGeom.StatorGeometry().OuterProfile().Diameter(stageGeom.StatorGeometry().XGapOut())
	var dOuterRotorIn = stageGeom.RotorGeometry().OuterProfile().Diameter(0)
	assert.True(
		suite.T(),
		common.ApproxEqual(dOuterStatorOut, dOuterRotorIn, 0.000001),
		testMessage(dOuterStatorOut, dOuterRotorIn),
	)
}


func TestStageGeometryGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(StageGeometryGeneratorTestSuite))
}
