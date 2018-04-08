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

func getTestCompressor() StagedCompressorNode {
	genFunc := func() IncompleteStageGeometryGenerator {
		return NewIncompleteStageGeomGen(
			NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
			NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut),
		)
	}
	gens := []IncompleteStageGeometryGenerator{genFunc(), genFunc(), genFunc()}

	node := NewStagedCompressorNode(
		rpm, dRelIn,
		true, gens,
		common.FromConst(htCoef),
		common.FromConst(reactivity),
		common.FromConst(labourCoef),
		common.FromConst(etaAd),
		common.FromConst(caCoef),
		precision, relaxCoef, 0.7, iterLimit,
	)

	node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node.TemperatureInput().SetState(states.NewTemperaturePortState(t0))
	node.PressureInput().SetState(states.NewPressurePortState(p0))
	node.MassRateInput().SetState(states.NewMassRatePortState(massRate))
	return node
}

type StagedCompressorNodeTestSuite struct {
	suite.Suite
	node *stagedCompressorNode
}

func (suite *StagedCompressorNodeTestSuite) SetupTest() {
	suite.node = getTestCompressor().(*stagedCompressorNode)
	suite.Require().Nil(suite.node.Process())
}

func (suite *StagedCompressorNodeTestSuite) TestSmoke() {
	//fmt.Println(PiStag(suite.node))
	msgs := make([]string, len(suite.node.Stages()))
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
