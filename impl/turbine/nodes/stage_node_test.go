package nodes

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/Sovianum/turbocycle/common"
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
)

type StageNodeTestSuite struct {
	suite.Suite
	pack *dataPack
	node *turbineStageNode
}

func (suite *StageNodeTestSuite) SetupTest() {
	suite.node = NewTurbineStageNode(n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision).(*turbineStageNode)
	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi / 2),
		states2.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	suite.node.PressureInput().SetState(states.NewPressurePortState(pg))

	suite.pack = new(dataPack)
}

func (suite *StageNodeTestSuite) TestT0() {
	var cp = getGas().Cp(tg)
	suite.node.t0(suite.pack)
	assert.True(suite.T(), common.ApproxEqual(tg - c0 * c0 / (2 * cp), suite.pack.T0, 0.01))
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}

func getGas() gases.Gas {
	return gases.GetAir()
}


