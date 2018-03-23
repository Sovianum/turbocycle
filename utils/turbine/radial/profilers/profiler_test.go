package profilers

import (
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/laws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	windage    = 1
	approxTRel = 0.7

	lRelOut  = 0.2
	bRel     = 3
	deltaRel = 0.1
	gammaIn  = -5
	gammaOut = 5

	cUIn  = 100
	cUOut = 110

	cAIn  = 100
	cAOut = 105

	u = 50

	installationAngle    = 45
	inletExpansionAngle  = 40
	outletExpansionAngle = 10

	inletFraction  = 0.5
	outletFraction = 0.5
)

type profilingFunc func(hRel float64) float64

type ProfilerTestSuite struct {
	suite.Suite
	behavior                 ProfilingBehavior
	geomGen                  geometry.BladingGeometryGenerator
	meanInletTriangle        states.VelocityTriangle
	meanOutletTriangle       states.VelocityTriangle
	inletVelocityLaw         laws.VelocityLaw
	outletVelocityLaw        laws.VelocityLaw
	profiler                 Profiler
	inletProfileAngleFunc    func(characteristicAngle, hRel float64) float64
	outletProfileAngleFunc   func(characteristicAngle, hRel float64) float64
	installationAngleFunc    profilingFunc
	inletExpansionAngleFunc  profilingFunc
	outletExpansionAngleFunc profilingFunc
	inletFractionFunc        profilingFunc
	outletFractionFunc       profilingFunc
}

func (suite *ProfilerTestSuite) SetupTest() {
	suite.behavior = NewStatorProfilingBehavior()
	suite.geomGen = geometry.NewGeneratorFromProfileAngles(
		lRelOut, bRel, deltaRel,
		common.ToRadians(gammaIn), common.ToRadians(gammaOut), approxTRel,
	)
	suite.meanInletTriangle = states.NewInletTriangleFromProjections(
		cUIn, cAIn, u,
	)
	suite.meanOutletTriangle = states.NewOutletTriangleFromProjections(
		cUOut, cAOut, u,
	)
	suite.inletVelocityLaw = laws.NewConstantCirculationVelocityLaw()
	suite.outletVelocityLaw = laws.NewConstantAbsoluteAngleLaw()
	suite.inletProfileAngleFunc = func(characteristicAngle, hRel float64) float64 {
		return characteristicAngle
	}
	suite.outletProfileAngleFunc = func(characteristicAngle, hRel float64) float64 {
		return characteristicAngle
	}
	suite.installationAngleFunc = func(hRel float64) float64 {
		return common.ToRadians(installationAngle)
	}
	suite.inletExpansionAngleFunc = func(hRel float64) float64 {
		return common.ToRadians(inletExpansionAngle)
	}
	suite.outletExpansionAngleFunc = func(hRel float64) float64 {
		return common.ToRadians(outletExpansionAngle)
	}
	suite.inletFractionFunc = func(hRel float64) float64 {
		return inletFraction
	}
	suite.outletFractionFunc = func(hRel float64) float64 {
		return outletFraction
	}

	suite.profiler = NewProfiler(
		windage, approxTRel,
		suite.behavior,
		suite.geomGen,
		suite.meanInletTriangle,
		suite.meanOutletTriangle,
		suite.inletVelocityLaw,
		suite.outletVelocityLaw,
		suite.inletProfileAngleFunc,
		suite.outletProfileAngleFunc,
		suite.installationAngleFunc,
		suite.inletExpansionAngleFunc,
		suite.outletExpansionAngleFunc,
		suite.inletFractionFunc,
		suite.outletFractionFunc,
	)
}

func (suite *ProfilerTestSuite) TestSmoke() {
	assert.NotNil(suite.T(), suite.profiler.InletTriangle(0))
	assert.NotNil(suite.T(), suite.profiler.OutletTriangle(0))
	suite.profiler.InletProfileAngle(0)
	suite.profiler.OutletProfileAngle(0)
	suite.profiler.InstallationAngle(0)
	suite.profiler.InletExpansionAngle(0)
	suite.profiler.OutletExpansionAngle(0)
	suite.profiler.InletPSAngleFraction(0)
	suite.profiler.OutletPSAngleFraction(0)
}

func TestProfilerTestSuite(t *testing.T) {
	suite.Run(t, new(ProfilerTestSuite))
}
