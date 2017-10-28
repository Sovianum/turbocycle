package nodes

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	n             = 1e4
	stageHeatDrop = 5e4
	reactivity    = 0.5
	phi           = 0.98
	psi           = 0.98
	airGapRel     = 0.001
	precision     = 0.05

	c0 = 400
	tg = 1200
	pg = 1e6

	gammaIn = -0.09
	gammaOut = 0.09
	baRel = 4
	lRelOut = 1 / 5
	deltaRel = 0.1
)

type StageNodeTestSuite struct {
	suite.Suite
	pack *dataPack
	node *turbineStageNode
}

func (suite *StageNodeTestSuite) SetupTest() {
	var generator = geometry.NewStageGeometryGenerator(
		geometry.NewGeneratorFromProfileAngles(lRelOut, baRel, deltaRel, gammaIn, gammaOut),
		geometry.NewGeneratorFromProfileAngles(lRelOut, baRel, deltaRel, gammaIn, gammaOut),
	)
	suite.node = NewTurbineStageNode(
		n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, generator,
	).(*turbineStageNode)
	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi / 2),
		states2.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	suite.node.PressureInput().SetState(states.NewPressurePortState(pg))

	suite.pack = new(dataPack)
}

func (suite *StageNodeTestSuite) TestGetStatorMeanInletDiameter() {
	//var density0 = 2.405

}

func (suite *StageNodeTestSuite) TestDensity0() {
	var t0 = 1130.
	var p0 = 7.8e5
	var density0 = p0 / (getGas().R() * t0)

	suite.pack.T0 = t0
	suite.pack.P0 = p0
	suite.node.density0(suite.pack)

	assert.InDelta(suite.T(), density0, suite.pack.Density0, 0.001)
}

func (suite *StageNodeTestSuite) TestP0() {
	var k = gases.K(getGas(), tg)
	suite.pack.T0 = 1130
	var p0 = pg * math.Pow(tg / suite.pack.T0, -k/(k-1))
	suite.node.p0(suite.pack)

	assert.InDelta(suite.T(), p0, suite.pack.P0, 0.01)
}

func (suite *StageNodeTestSuite) TestT0() {
	var cp = getGas().Cp(tg)
	suite.node.t0(suite.pack)
	assert.InDelta(suite.T(), tg - c0 * c0 / (2 * cp), suite.pack.T0, 0.01)
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}

func getGas() gases.Gas {
	return gases.GetAir()
}


