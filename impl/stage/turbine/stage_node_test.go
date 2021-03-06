package turbine

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	states2 "github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	n             = 1e4
	stageHeatDrop = 2.5e5
	reactivity    = 0.5
	phi           = 0.98
	psi           = 0.985
	airGapRel     = 0.001
	precision     = 1e-4

	c0       = 50.
	tg       = 1200.
	pg       = 1e6
	massRate = 100.
	baRel    = 4

	t0       = 1130.
	p0       = 7.8e5
	density0 = 2.405
)

type StageNodeTestSuite struct {
	suite.Suite
	pack *DataPack
	node *turbineStageNode
	gen  StageGeometryGenerator
}

func (suite *StageNodeTestSuite) SetupTest() {
	suite.gen = NewStageGeometryGenerator(
		lRelOut,
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut, statorApproxTRel),
		NewIncompleteGenerator(baRel, deltaRel, gammaIn, gammaOut, rotorApproxTRel),
	)

	suite.node = NewTurbineSingleStageNode(
		n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, suite.gen,
	).(*turbineStageNode)
	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))

	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	suite.node.PressureInput().SetState(states.NewPressurePortState(pg))
	suite.node.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	suite.pack = new(DataPack)
}

func (suite *StageNodeTestSuite) TestCalcFirstStage() {
	suite.node.SetStageMode(first)
	suite.node.SetAlpha1FirstStage(common.ToRadians(14))
	var pack = suite.node.getDataPack()

	suite.False(math.IsNaN(pack.U2))
	suite.False(math.IsNaN(pack.Tw1))
	suite.False(math.IsNaN(pack.Pw1))
	suite.False(math.IsNaN(pack.RotorHeatDrop))
	suite.False(math.IsNaN(pack.WAd2))
	suite.False(math.IsNaN(pack.W2))
	suite.False(math.IsNaN(pack.T2))
	suite.False(math.IsNaN(pack.T2Prime))
	suite.False(math.IsNaN(pack.P2))
	suite.False(math.IsNaN(pack.Density2))
	suite.False(math.IsNaN(pack.C2a))
	suite.False(math.IsNaN(pack.Beta2))
	suite.False(math.IsNaN(pack.C2u))
	suite.False(math.IsNaN(pack.Pi))
	suite.False(math.IsNaN(pack.MeanRadiusLabour))
	suite.False(math.IsNaN(pack.EtaU))
	suite.False(math.IsNaN(pack.StatorSpecificLoss))
	suite.False(math.IsNaN(pack.RotorSpecificLoss))
	suite.False(math.IsNaN(pack.OutletVelocitySpecificLoss))
	suite.False(math.IsNaN(pack.AirGapSpecificLoss))
	suite.False(math.IsNaN(pack.VentilationSpecificLoss))
	suite.False(math.IsNaN(pack.EtaT))
	suite.False(math.IsNaN(pack.T2Stag))
	suite.False(math.IsNaN(pack.P2Stag))
	suite.False(math.IsNaN(pack.StageLabour))
	suite.False(math.IsNaN(pack.StageHeatDropStag))
	suite.False(math.IsNaN(pack.EtaT))
}

func (suite *StageNodeTestSuite) TestInitCalcFirstStage() {
	suite.node.SetStageMode(first)
	suite.node.SetAlpha1FirstStage(common.ToRadians(14))
	suite.node.initCalcFirst(suite.pack)
	suite.False(math.IsNaN(suite.pack.RotorHeatDrop))
	suite.False(math.IsNaN(suite.pack.T1Prime))
	suite.False(math.IsNaN(suite.pack.C1Ad))
	suite.False(math.IsNaN(suite.pack.C1))
	suite.False(math.IsNaN(suite.pack.C1a))
	suite.False(math.IsNaN(suite.pack.T1))
	suite.False(math.IsNaN(suite.pack.P1))
	suite.False(math.IsNaN(suite.pack.Density1))
	suite.False(math.IsNaN(suite.pack.Area1))
	suite.False(math.IsNaN(suite.pack.RotorMeanInletDiameter))
	suite.False(math.IsNaN(suite.pack.U1))
}

func (suite *StageNodeTestSuite) TestInitCalc() {
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, 10, math.Pi/2),
		states2.InletTriangleType,
	))
	suite.node.initCalcOnly(suite.pack)

	suite.False(math.IsNaN(suite.pack.T0))
	suite.False(math.IsNaN(suite.pack.P0))
	suite.False(math.IsNaN(suite.pack.Density0))
	suite.False(math.IsNaN(suite.pack.StatorMeanInletDiameter))
	suite.False(math.IsNaN(suite.pack.StatorHeatDrop))
	suite.False(math.IsNaN(suite.pack.T1Prime))
	suite.False(math.IsNaN(suite.pack.C1Ad))
	suite.False(math.IsNaN(suite.pack.C1))
	suite.False(math.IsNaN(suite.pack.T1))
	suite.False(math.IsNaN(suite.pack.P1))
	suite.False(math.IsNaN(suite.pack.Density1))

	suite.False(math.IsNaN(suite.pack.Area1))
	suite.False(math.IsNaN(suite.pack.C1a))
	suite.False(math.IsNaN(suite.pack.U1))
	suite.False(math.IsNaN(suite.pack.Alpha1))
}

func (suite *StageNodeTestSuite) TestStatorHeatDrop() {
	suite.node.thermo1(suite.pack)
	suite.True(
		common.ApproxEqual(stageHeatDrop*(1-reactivity), suite.pack.StatorHeatDrop, 0.00001),
	)
}

func (suite *StageNodeTestSuite) TestGetStatorMeanInletDiameter() {
	suite.pack.Density0 = density0

	baRel := suite.gen.StatorGenerator().Elongation()
	_, gammaMean := geometry.GetTotalAndMeanLineAngles(
		gammaIn, gammaOut, MidLineFactor,
	)

	enom := baRel - (1+deltaRel)*(math.Tan(gammaOut)-math.Tan(gammaIn))
	denom := baRel - 2*(1+deltaRel)*lRelOut*math.Tan(gammaMean)
	lRelInRotor := lRelOut * enom / denom

	lRelInStator := RecalculateLRel(lRelInRotor, -(1+deltaRel)/elongation, gammaIn, gammaOut)
	expectedDMean := math.Sqrt(massRate / (math.Pi * c0 * density0 * lRelInStator))

	suite.node.getStageGeometry(suite.pack)

	suite.Require().Equal(c0, suite.node.statorInletTriangle().C())
	suite.InDelta(expectedDMean, suite.pack.StatorMeanInletDiameter, 1e-4)
}

func (suite *StageNodeTestSuite) TestDensity0() {
	var density0 = p0 / (getGas().R() * t0)
	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(t0))
	suite.node.PressureInput().SetState(states.NewPressurePortState(p0))
	suite.node.thermo0(suite.pack)

	assert.InDelta(suite.T(), density0, suite.pack.Density0, 0.1) // low precision cos use static parameters
}

func (suite *StageNodeTestSuite) TestT0() {
	var cp = getGas().Cp(tg)
	suite.node.thermo0(suite.pack)
	assert.True(
		suite.T(),
		common.ApproxEqual(tg-c0*c0/(2*cp), suite.pack.T0, 0.0001),
	)
}

func (suite *StageNodeTestSuite) TestAlpha1Trend() {
	alpha11 := common.ToRadians(14)
	alpha12 := common.ToRadians(16)

	node1 := NewTurbineFirstStageNode(alpha11, n, stageHeatDrop, reactivity, phi, psi, airGapRel, 1e-5, suite.gen)
	node2 := NewTurbineFirstStageNode(alpha12, n, stageHeatDrop, reactivity, phi, psi, airGapRel, 1e-5, suite.gen)

	node1.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	node1.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node1.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	node1.PressureInput().SetState(states.NewPressurePortState(pg))
	node1.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	node2.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	node2.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node2.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	node2.PressureInput().SetState(states.NewPressurePortState(pg))
	node2.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	suite.Require().NoError(node1.Process())
	suite.Require().NoError(node2.Process())

	suite.True(
		node1.GetDataPack().RotorOutletTriangle.Alpha() > node2.GetDataPack().RotorOutletTriangle.Alpha(),
		"alpha_21 = %.3f, alpha_22 = %.3f",
		common.ToDegrees(node1.GetDataPack().RotorOutletTriangle.Alpha()),
		common.ToDegrees(node2.GetDataPack().RotorOutletTriangle.Alpha()),
	)
}

func (suite *StageNodeTestSuite) TestReactivityTrend() {
	reactivity1 := 0.5
	reactivity2 := 0.6
	dMean := 0.75

	node1 := NewTurbineMidStageNode(dMean, n, stageHeatDrop, reactivity1, phi, psi, airGapRel, 1e-5, suite.gen)
	node2 := NewTurbineMidStageNode(dMean, n, stageHeatDrop, reactivity2, phi, psi, airGapRel, 1e-5, suite.gen)

	node1.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	node1.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node1.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	node1.PressureInput().SetState(states.NewPressurePortState(pg))
	node1.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	node2.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	node2.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))
	node2.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	node2.PressureInput().SetState(states.NewPressurePortState(pg))
	node2.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	suite.Require().NoError(node1.Process())
	suite.Require().NoError(node2.Process())

	suite.True(
		node1.GetDataPack().RotorInletTriangle.Alpha() < node2.GetDataPack().RotorInletTriangle.Alpha(),
		"alpha_11 = %.3f, alpha_12 = %.3f",
		common.ToDegrees(node1.GetDataPack().RotorInletTriangle.Alpha()),
		common.ToDegrees(node2.GetDataPack().RotorInletTriangle.Alpha()),
	)
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageNodeTestSuite))
}

func getGas() gases.Gas {
	return gases.GetAir()
}
