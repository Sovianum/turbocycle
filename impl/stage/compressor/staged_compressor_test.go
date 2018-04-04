package compressor

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	states2 "github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	stageNum = 3
)

type StagedCompressorNodeTestSuite struct {
	suite.Suite
	node *stagedCompressorNode
	gens []StageGeometryGenerator
}

func (suite *StagedCompressorNodeTestSuite) SetupTest() {
	genFunc := func() StageGeometryGenerator {
		return NewStageGeometryGenerator(
			dRelIn,
			NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
			NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
		)
	}
	suite.gens = []StageGeometryGenerator{genFunc(), genFunc(), genFunc()}

	suite.node = NewStagedCompressorNode(
		rpm, dRelIn,
		suite.gens,
		common.FromConst(htCoef),
		common.FromConst(reactivity),
		common.FromConst(labourCoef),
		common.FromConst(etaAd),
		common.FromConst(caCoef),
		precision, relaxCoef, 0.7, iterLimit,
	).(*stagedCompressorNode)

	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(t0))
	suite.node.PressureInput().SetState(states.NewPressurePortState(p0))
	suite.node.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	suite.Require().Nil(suite.node.Process())
}

func (suite *StagedCompressorNodeTestSuite) TestSmoke() {
	msgs := make([]string, stageNum)
	for i, stage := range suite.node.Stages() {
		b, e := json.MarshalIndent(stage.GetDataPack(), "", "\t")
		suite.Require().NoError(e)
		msgs[i] = string(b)
	}
	fmt.Println(msgs)
}

func TestStagedCompressorNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StagedCompressorNodeTestSuite))
}
