package compressor

import (
	"math"
	"testing"

	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	rpm        = 1e4
	reactivity = 0.5

	c0       = 50.
	massRate = 100.
	baRel    = 4

	t0 = 288
	p0 = 1e5

	deltaRel = 0.1

	gammaOut = 0
	gammaIn  = 0.5

	dRelIn     = 0.5
	htCoef     = 0.2
	labourCoef = 1
	etaAd      = 0.82
	caCoef     = 0.5

	precision = 1e-5
	iterLimit = 100
	relaxCoef = 1
)

type StageNodeTestSuite struct {
	suite.Suite
	pack *DataPack
	node *stageNode
	gen  StageGeometryGenerator
}

func (suite *StageNodeTestSuite) SetupTest() {
	suite.initNode(rpm)
}

func (suite *StageNodeTestSuite) initNode(rpmArg float64) {
	suite.gen = NewStageGeometryGenerator(
		dRelIn,
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
	)

	suite.node = NewFirstStageNode(
		dRelIn,
		htCoef, htCoef,
		reactivity, reactivity,
		labourCoef, etaAd, caCoef, rpmArg,
		suite.gen, precision, relaxCoef, 0.99, iterLimit,
	).(*stageNode)
	suite.node.GasInput().SetState(states2.NewGasPortState(gases.GetAir()))

	suite.node.VelocityInput().SetState(states.NewVelocityPortState(
		states.NewInletTriangle(0, c0, math.Pi/2),
		states.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states2.NewTemperaturePortState(t0))
	suite.node.PressureInput().SetState(states2.NewPressurePortState(p0))
	suite.node.MassRateInput().SetState(states2.NewMassRatePortState(massRate))

	suite.Require().Nil(suite.node.Process())
	suite.pack = suite.node.pack
}

func (suite *StageNodeTestSuite) TestGeom() {
	rotorGeom := suite.node.pack.StageGeometry.RotorGeometry()
	rotorWidth := rotorGeom.XGapOut()

	statorGeom := suite.node.pack.StageGeometry.StatorGeometry()
	statorWidth := statorGeom.XGapOut()

	suite.EqualValues(rotorGeom.InnerProfile().Diameter(rotorWidth), statorGeom.InnerProfile().Diameter(0))

	stageWidth := rotorWidth + statorWidth
	dInIn := rotorGeom.InnerProfile().Diameter(0)

	suite.InDelta(dInIn+2*rotorWidth*math.Tan(gammaIn), rotorGeom.InnerProfile().Diameter(rotorWidth), 1e-9)
	suite.InDelta(dInIn+2*stageWidth*math.Tan(gammaIn), statorGeom.InnerProfile().Diameter(statorWidth), 1e-9)
}

func (suite *StageNodeTestSuite) TestSmoke() {
	//if b, e := json.MarshalIndent(suite.pack, "", "\t"); e != nil {
	//	suite.Require().Nil(e)
	//} else {
	//	fmt.Println(string(b))
	//}
}

func (suite *StageNodeTestSuite) TestTriangles() {
	suite.initNode(rpm)
	triangle11 := suite.pack.InletTriangle
	triangle12 := suite.pack.OutletTriangle
	geom1 := suite.pack.StageGeometry.StatorGeometry()
	r13 := geometry.RRel(
		geom1.InnerProfile().Diameter(geom1.XGapOut()) / geom1.OuterProfile().Diameter(geom1.XGapOut()),
	)

	suite.initNode(rpm * 0.6)
	triangle21 := suite.pack.InletTriangle
	triangle22 := suite.pack.OutletTriangle
	geom2 := suite.pack.StageGeometry.StatorGeometry()
	r23 := geometry.RRel(
		geom2.InnerProfile().Diameter(geom2.XGapOut()) / geom2.OuterProfile().Diameter(geom2.XGapOut()),
	)

	suite.InDelta(triangle11.Alpha(), triangle21.Alpha(), 1e-9)
	suite.InDelta(triangle11.Beta(), triangle21.Beta(), 1e-9)

	suite.True(triangle12.CA() > triangle22.CA())
	suite.InDelta(r13, r23, 1e-9)

	suite.True(triangle12.Alpha() < triangle22.Alpha(), "%f, %f", triangle12.Alpha(), triangle22.Alpha())
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}
