package turbine

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	common2 "github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	states2 "github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
)

const (
	alpha1 = 20. / 180. * math.Pi
)

func getTestTurbine() StagedTurbineNode {
	genFunc := func() IncompleteStageGeometryGenerator {
		return NewIncompleteStageGeometryGenerator(
			NewIncompleteGenerator(elongation, deltaRel, gammaIn, gammaOut, approxTRel),
			NewIncompleteGenerator(elongation, deltaRel, gammaIn, gammaOut, approxTRel),
		)
	}
	gens := []IncompleteStageGeometryGenerator{genFunc(), genFunc(), genFunc()}

	node := NewStagedTurbineNode(
		n, alpha1, stageHeatDrop*2, 0.2,
		common.Func1DFromConst(phi),
		common.Func1DFromConst(psi),
		common.Func1DFromConst(reactivity),
		common.Func1DFromConst(airGapRel),
		common.Func1DFromConst(1),
		gens, precision,
	)
	node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))

	node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	node.PressureInput().SetState(states.NewPressurePortState(pg))
	node.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	return node
}

type StagedTurbineTestSuite struct {
	suite.Suite
	node *stagedTurbineNode
}

func (suite *StagedTurbineTestSuite) SetupTest() {
	suite.node = getTestTurbine().(*stagedTurbineNode)
	suite.Require().NoError(suite.node.Process())
}

func (suite *StagedTurbineTestSuite) TestSmoke() {
	msgs := make([]string, 3)
	for i, stage := range suite.node.Stages() {
		b, e := json.MarshalIndent(stage.GetDataPack(), "", "\t")
		suite.Require().NoError(e)
		msgs[i] = string(b)
	}

	for i, stage := range suite.node.stages {
		geom := stage.GetDataPack().StageGeometry
		stator := geom.StatorGeometry()
		rotor := geom.RotorGeometry()
		fmt.Printf(
			"#%d: d_in_in: %.3f, d_out_in: %.3f, d_in_out: %.3f, d_out_out: %.3f\n",
			i,
			stator.InnerProfile().Diameter(0),
			stator.OuterProfile().Diameter(0),
			rotor.InnerProfile().Diameter(rotor.XGapOut()),
			rotor.OuterProfile().Diameter(rotor.XGapOut()),
		)
	}
	fmt.Println()
	for i, stage := range suite.node.Stages() {
		fmt.Printf(
			"#%d: inlet: %s\t outlet: %s\n",
			i,
			SprintfTriangle(stage.GetDataPack().RotorInletTriangle),
			SprintfTriangle(stage.GetDataPack().RotorOutletTriangle),
		)
	}
}

func TestStagedTurbineTestSuite(t *testing.T) {
	suite.Run(t, new(StagedTurbineTestSuite))
}

func SprintfTriangle(triangle states2.VelocityTriangle) string {
	return fmt.Sprintf(
		"alpha: %.3f; beta: %.3f",
		common2.ToDegrees(triangle.Alpha()),
		common2.ToDegrees(triangle.Beta()),
	)
}
