package compressor

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
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
	gammaIn  = 0.06

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
	suite.gen = NewStageGeometryGenerator(
		dRelIn,
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
	)

	suite.node = NewStageNode(
		dRelIn,
		htCoef, htCoef,
		reactivity, reactivity,
		labourCoef, etaAd, caCoef, rpm,
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

func (suite *StageNodeTestSuite) TestSmoke() {
	if b, e := json.MarshalIndent(suite.pack, "", "\t"); e != nil {
		suite.Require().Nil(e)
	} else {
		fmt.Println(string(b))
	}
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}
