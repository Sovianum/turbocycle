package nodes

import (
	"math"
	"testing"

	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	n             = 1e4
	stageHeatDrop = 3e5
	reactivity    = 0.5
	phi           = 0.98
	psi           = 0.98
	airGapRel     = 0.001
	precision     = 0.05

	c0       = 50.
	tg       = 1200.
	pg       = 1e6
	massRate = 100.

	gammaIn  = -0.09
	gammaOut = 0.09
	baRel    = 4
	lRelOut  = 0.15
	deltaRel = 0.1

	t0       = 1130.
	p0       = 7.8e5
	density0 = 2.405

	statorApproxTRel = 0.7
	rotorApproxTRel = 0.7

	alpha = 0.5
)

type StageNodeTestSuite struct {
	suite.Suite
	pack *DataPack
	node *turbineStageNode
	gen  geometry.StageGeometryGenerator
}

func (suite *StageNodeTestSuite) SetupTest() {
	suite.gen = geometry.NewStageGeometryGenerator(
		lRelOut,
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, statorApproxTRel),
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, rotorApproxTRel),
	)

	suite.node = NewTurbineStageNode(
		n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, suite.gen,
	).(*turbineStageNode)
	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))

	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	suite.node.PressureInput().SetState(states.NewPressurePortState(pg))
	suite.node.MassRateInput().SetState(states2.NewMassRatePortState(massRate))

	suite.pack = new(DataPack)
}

func (suite *StageNodeTestSuite) TestCalc() {
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, 10, math.Pi/2),
		states2.InletTriangleType,
	))
	var pack = suite.node.getDataPack()

	assert.False(suite.T(), math.IsNaN(pack.U2))
	assert.False(suite.T(), math.IsNaN(pack.Tw1))
	assert.False(suite.T(), math.IsNaN(pack.Pw1))
	assert.False(suite.T(), math.IsNaN(pack.RotorHeatDrop))
	assert.False(suite.T(), math.IsNaN(pack.WAd2))
	assert.False(suite.T(), math.IsNaN(pack.W2))
	assert.False(suite.T(), math.IsNaN(pack.T2))
	assert.False(suite.T(), math.IsNaN(pack.T2Prime))
	assert.False(suite.T(), math.IsNaN(pack.P2))
	assert.False(suite.T(), math.IsNaN(pack.Density2))
	assert.False(suite.T(), math.IsNaN(pack.C2a))
	assert.False(suite.T(), math.IsNaN(pack.Beta2))
	assert.False(suite.T(), math.IsNaN(pack.C2u))
	assert.False(suite.T(), math.IsNaN(pack.Pi))
	assert.False(suite.T(), math.IsNaN(pack.MeanRadiusLabour))
	assert.False(suite.T(), math.IsNaN(pack.EtaU))
	assert.False(suite.T(), math.IsNaN(pack.StatorSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.RotorSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.OutletVelocitySpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.AirGapSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.VentilationSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.EtaT))
	assert.False(suite.T(), math.IsNaN(pack.T2Stag))
	assert.False(suite.T(), math.IsNaN(pack.P2Stag))
	assert.False(suite.T(), math.IsNaN(pack.StageLabour))
	assert.False(suite.T(), math.IsNaN(pack.StageHeatDropStag))
	assert.False(suite.T(), math.IsNaN(pack.EtaT))
}

func (suite *StageNodeTestSuite) TestCalcFirstStage() {
	suite.node.SetFirstStageMode(true)
	suite.node.SetAlpha1FirstStage(common.ToRadians(14))
	var pack = suite.node.getDataPack()

	assert.False(suite.T(), math.IsNaN(pack.U2))
	assert.False(suite.T(), math.IsNaN(pack.Tw1))
	assert.False(suite.T(), math.IsNaN(pack.Pw1))
	assert.False(suite.T(), math.IsNaN(pack.RotorHeatDrop))
	assert.False(suite.T(), math.IsNaN(pack.WAd2))
	assert.False(suite.T(), math.IsNaN(pack.W2))
	assert.False(suite.T(), math.IsNaN(pack.T2))
	assert.False(suite.T(), math.IsNaN(pack.T2Prime))
	assert.False(suite.T(), math.IsNaN(pack.P2))
	assert.False(suite.T(), math.IsNaN(pack.Density2))
	assert.False(suite.T(), math.IsNaN(pack.C2a))
	assert.False(suite.T(), math.IsNaN(pack.Beta2))
	assert.False(suite.T(), math.IsNaN(pack.C2u))
	assert.False(suite.T(), math.IsNaN(pack.Pi))
	assert.False(suite.T(), math.IsNaN(pack.MeanRadiusLabour))
	assert.False(suite.T(), math.IsNaN(pack.EtaU))
	assert.False(suite.T(), math.IsNaN(pack.StatorSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.RotorSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.OutletVelocitySpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.AirGapSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.VentilationSpecificLoss))
	assert.False(suite.T(), math.IsNaN(pack.EtaT))
	assert.False(suite.T(), math.IsNaN(pack.T2Stag))
	assert.False(suite.T(), math.IsNaN(pack.P2Stag))
	assert.False(suite.T(), math.IsNaN(pack.StageLabour))
	assert.False(suite.T(), math.IsNaN(pack.StageHeatDropStag))
	assert.False(suite.T(), math.IsNaN(pack.EtaT))
}

func (suite *StageNodeTestSuite) TestInitCalcFirstStage() {
	suite.node.SetFirstStageMode(true)
	suite.node.SetAlpha1FirstStage(common.ToRadians(14))
	suite.node.initCalcFirstStage(suite.pack)
	assert.False(suite.T(), math.IsNaN(suite.pack.RotorHeatDrop))
	assert.False(suite.T(), math.IsNaN(suite.pack.T1Prime))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1Ad))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1a))
	assert.False(suite.T(), math.IsNaN(suite.pack.T1))
	assert.False(suite.T(), math.IsNaN(suite.pack.P1))
	assert.False(suite.T(), math.IsNaN(suite.pack.Density1))
	assert.False(suite.T(), math.IsNaN(suite.pack.Area1))
	assert.False(suite.T(), math.IsNaN(suite.pack.RotorMeanInletDiameter))
	assert.False(suite.T(), math.IsNaN(suite.pack.U1))
}

func (suite *StageNodeTestSuite) TestInitCalc() {
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, 10, math.Pi/2),
		states2.InletTriangleType,
	))
	suite.node.initCalc(suite.pack)

	assert.False(suite.T(), math.IsNaN(suite.pack.T0))
	assert.False(suite.T(), math.IsNaN(suite.pack.P0))
	assert.False(suite.T(), math.IsNaN(suite.pack.Density0))
	assert.False(suite.T(), math.IsNaN(suite.pack.StatorMeanInletDiameter))
	assert.False(suite.T(), math.IsNaN(suite.pack.StatorHeatDrop))
	assert.False(suite.T(), math.IsNaN(suite.pack.T1Prime))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1Ad))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1))
	assert.False(suite.T(), math.IsNaN(suite.pack.T1))
	assert.False(suite.T(), math.IsNaN(suite.pack.P1))
	assert.False(suite.T(), math.IsNaN(suite.pack.Density1))

	assert.False(suite.T(), math.IsNaN(suite.pack.Area1))
	assert.False(suite.T(), math.IsNaN(suite.pack.C1a))
	assert.False(suite.T(), math.IsNaN(suite.pack.U1))
	assert.False(suite.T(), math.IsNaN(suite.pack.Alpha1))
}

func (suite *StageNodeTestSuite) TestStatorHeatDrop() {
	suite.node.statorHeatDrop(suite.pack)
	assert.True(
		suite.T(),
		common.ApproxEqual(stageHeatDrop*(1-reactivity), suite.pack.StatorHeatDrop, 0.00001),
	)
}

func (suite *StageNodeTestSuite) TestGetStatorMeanInletDiameter() {
	suite.pack.Density0 = density0

	var baRel = suite.gen.StatorGenerator().Elongation()
	var _, gammaMean = geometry.GetTotalAndMeanLineAngles(
		gammaIn, gammaOut,
	)

	var enom = baRel - (1+deltaRel)*(math.Tan(gammaOut)-math.Tan(gammaIn))
	var denom = baRel - 2*(1+deltaRel)*lRelOut*math.Tan(gammaMean)
	var lRelIn = enom / denom
	var expectedDMean = math.Sqrt(massRate / (math.Pi * c0 * density0 * lRelIn))

	suite.node.getStatorMeanInletDiameter(suite.pack)

	assert.True(
		suite.T(),
		common.ApproxEqual(expectedDMean, suite.pack.StatorMeanInletDiameter, 0.00001),
	)
}

func (suite *StageNodeTestSuite) TestDensity0() {
	var density0 = p0 / (getGas().R() * t0)

	suite.pack.T0 = t0
	suite.pack.P0 = p0
	suite.node.density0(suite.pack)

	assert.InDelta(suite.T(), density0, suite.pack.Density0, 0.001)
}

func (suite *StageNodeTestSuite) TestP0() {
	var k = gases.K(getGas(), tg)
	var expectedP0 = pg * math.Pow(tg/t0, -k/(k-1))

	suite.pack.T0 = t0
	suite.node.p0(suite.pack)

	assert.True(
		suite.T(),
		common.ApproxEqual(expectedP0, suite.pack.P0, 0.0001),
		testMessage(expectedP0, suite.pack.P0),
	)
}

func (suite *StageNodeTestSuite) TestT0() {
	var cp = getGas().Cp(tg)
	suite.node.t0(suite.pack)
	assert.True(
		suite.T(),
		common.ApproxEqual(tg-c0*c0/(2*cp), suite.pack.T0, 0.0001),
	)
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}

func getGas() gases.Gas {
	return gases.GetAir()
}

func testMessage(x0, x float64) string {
	return fmt.Sprintf("expected: %v; got %v", x0, x)
}
