package profilers

import (
	"math"

	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/laws"
)

func NewProfiler(
	windage float64,
	approxTRel float64,

	behavior ProfilingBehavior,
	geomGen turbine.BladingGeometryGenerator,

	meanInletTriangle states.VelocityTriangle,
	meanOutletTriangle states.VelocityTriangle,

	inletVelocityLaw laws.InletVelocityLaw,
	outletVelocityLaw laws.OutletVelocityLaw,

	inletProfileAngleFunc func(characteristicAngle, hRel float64) float64,
	outletProfileAngleFunc func(characteristicAngle, hRel float64) float64,

	installationAngleFunc func(hRel float64) float64,

	inletExpansionAngleFunc func(hRel float64) float64,
	outletExpansionAngleFunc func(hRel float64) float64,

	inletPSAngleFractionFunc func(hRel float64) float64,
	outletPSAngleFractionFunc func(hRel float64) float64,
) Profiler {
	return &profiler{
		windage:    windage,
		approxTRel: approxTRel,

		behavior: behavior,
		geomGen:  geomGen,

		meanInletTriangle:  meanInletTriangle,
		meanOutletTriangle: meanOutletTriangle,

		inletVelocityLaw:  inletVelocityLaw,
		outletVelocityLaw: outletVelocityLaw,

		inletProfileAngleFunc:  inletProfileAngleFunc,
		outletProfileAngleFunc: outletProfileAngleFunc,

		installationAngleFunc: installationAngleFunc,

		inletExpansionAngleFunc:  inletExpansionAngleFunc,
		outletExpansionAngleFunc: outletExpansionAngleFunc,

		inletPSAngleFractionFunc:  inletPSAngleFractionFunc,
		outletPSAngleFractionFunc: outletPSAngleFractionFunc,
	}
}

type Profiler interface {
	InletTriangle(hRel float64) states.VelocityTriangle
	OutletTriangle(hRel float64) states.VelocityTriangle

	InletProfileAngle(hRel float64) float64
	OutletProfileAngle(hRel float64) float64

	InstallationAngle(hRel float64) float64

	InletExpansionAngle(hRel float64) float64
	OutletExpansionAngle(hRel float64) float64

	InletPSAngleFraction(hRel float64) float64
	OutletPSAngleFraction(hRel float64) float64
}

func InletBendAngle(hRel float64, profiler Profiler) float64 {
	var installationAngle = profiler.InstallationAngle(hRel)
	var inletProfileAngle = profiler.InletProfileAngle(hRel)
	return math.Pi - (installationAngle + inletProfileAngle)
}

func OutletBendAngle(hRel float64, profiler Profiler) float64 {
	var installationAngle = profiler.InstallationAngle(hRel)
	var outletProfileAngle = profiler.OutletProfileAngle(hRel)
	return installationAngle - outletProfileAngle
}

func InletPSExpansionAngle(hRel float64, profiler Profiler) float64 {
	var psAngleFraction = profiler.InletPSAngleFraction(hRel)
	var inletExpansionAngle = profiler.InletExpansionAngle(hRel)
	return psAngleFraction * inletExpansionAngle
}

func InletSSExpansionAngle(hRel float64, profiler Profiler) float64 {
	var psAngleFraction = profiler.InletPSAngleFraction(hRel)
	var inletExpansionAngle = profiler.InletExpansionAngle(hRel)
	return (1 - psAngleFraction) * inletExpansionAngle
}

func OutletPSExpansionAngle(hRel float64, profiler Profiler) float64 {
	var psAngleFraction = profiler.OutletPSAngleFraction(hRel)
	var expansionAngle = profiler.OutletExpansionAngle(hRel)
	return psAngleFraction * expansionAngle
}

func OutletSSExpansionAngle(hRel float64, profiler Profiler) float64 {
	var psAngleFraction = profiler.OutletPSAngleFraction(hRel)
	var expansionAngle = profiler.OutletExpansionAngle(hRel)
	return (1 - psAngleFraction) * expansionAngle
}

func InletPSAngle(hRel float64, profiler Profiler) float64 {
	var bendAngle = InletBendAngle(hRel, profiler)
	var psExpansionAngle = InletPSExpansionAngle(hRel, profiler)
	return bendAngle - psExpansionAngle
}

func InletSSAngle(hRel float64, profiler Profiler) float64 {
	var bendAngle = InletBendAngle(hRel, profiler)
	var ssExpansionAngle = InletSSExpansionAngle(hRel, profiler)
	return bendAngle + ssExpansionAngle
}

func OutletPSAngle(hRel float64, profiler Profiler) float64 {
	var bendAngle = OutletBendAngle(hRel, profiler)
	var psExpansionAngle = OutletPSExpansionAngle(hRel, profiler)
	return bendAngle - psExpansionAngle
}

func OutletSSAngle(hRel float64, profiler Profiler) float64 {
	return OutletBendAngle(hRel, profiler) + OutletSSExpansionAngle(hRel, profiler)
}

type profiler struct {
	behavior ProfilingBehavior

	geomGen turbine.BladingGeometryGenerator

	meanInletTriangle  states.VelocityTriangle
	meanOutletTriangle states.VelocityTriangle

	inletVelocityLaw  laws.InletVelocityLaw
	outletVelocityLaw laws.OutletVelocityLaw

	inletProfileAngleFunc  func(characteristicAngle, hRel float64) float64
	outletProfileAngleFunc func(characteristicAngle, hRel float64) float64

	installationAngleFunc func(hRel float64) float64

	inletExpansionAngleFunc  func(hRel float64) float64
	outletExpansionAngleFunc func(hRel float64) float64

	inletPSAngleFractionFunc  func(hRel float64) float64
	outletPSAngleFractionFunc func(hRel float64) float64

	windage    float64
	approxTRel float64
}

func (profiler *profiler) InletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.inletTriangle(hRel)
}

func (profiler *profiler) OutletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.outletTriangle(hRel)
}

func (profiler *profiler) InletProfileAngle(hRel float64) float64 {
	return profiler.inletProfileAngleFunc(
		profiler.behavior.ProfilingAngle(profiler.inletTriangle(hRel)),
		hRel,
	)
}

func (profiler *profiler) OutletProfileAngle(hRel float64) float64 {
	return profiler.outletProfileAngleFunc(
		profiler.behavior.ProfilingAngle(profiler.outletTriangle(hRel)),
		hRel,
	)
}

func (profiler *profiler) InstallationAngle(hRel float64) float64 {
	return profiler.installationAngleFunc(hRel)
}

func (profiler *profiler) InletExpansionAngle(hRel float64) float64 {
	return profiler.inletExpansionAngleFunc(hRel)
}

func (profiler *profiler) OutletExpansionAngle(hRel float64) float64 {
	return profiler.outletExpansionAngleFunc(hRel)
}

func (profiler *profiler) InletPSAngleFraction(hRel float64) float64 {
	return profiler.inletPSAngleFractionFunc(hRel)
}

func (profiler *profiler) OutletPSAngleFraction(hRel float64) float64 {
	return profiler.outletPSAngleFractionFunc(hRel)
}

func (profiler *profiler) inletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.inletVelocityLaw.InletTriangle(
		profiler.meanInletTriangle,
		hRel,
		turbine.LRelIn(profiler.geomGen),
	)
}

func (profiler *profiler) outletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.outletVelocityLaw.OutletTriangle(
		profiler.meanOutletTriangle,
		hRel,
		profiler.geomGen.LRelOut(),
	)
}
