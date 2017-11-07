package profilers

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/radial"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
)

func NewProfiler(
	windage float64,
	approxTRel float64,

	behavior ProfilingBehavior,
	geomGen geometry.BladingGeometryGenerator,

	meanInletTriangle states.VelocityTriangle,
	meanOutletTriangle states.VelocityTriangle,

	inletVelocityLaw radial.VelocityLaw,
	outletVelocityLaw radial.VelocityLaw,

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

	BladeNumber() int
}

func InletBendAngle(hRel float64, profiler Profiler) float64 {
	return math.Pi - (profiler.InstallationAngle(hRel) + profiler.InletProfileAngle(hRel))
}

func OutletBendAngle(hRel float64, profiler Profiler) float64 {
	return profiler.InstallationAngle(hRel) - profiler.OutletProfileAngle(hRel)
}

func InletPSExpansionAngle(hRel float64, profiler Profiler) float64 {
	return profiler.InletPSAngleFraction(hRel) * profiler.InletExpansionAngle(hRel)
}

func InletSSExpansionAngle(hRel float64, profiler Profiler) float64 {
	return (1 - profiler.InletPSAngleFraction(hRel)) * profiler.InletExpansionAngle(hRel)
}

func OutletPSExpansionAngle(hRel float64, profiler Profiler) float64 {
	return profiler.OutletPSAngleFraction(hRel) * profiler.InletExpansionAngle(hRel)
}

func OutletSSExpansionAngle(hRel float64, profiler Profiler) float64 {
	return (1 - profiler.OutletPSAngleFraction(hRel)) * profiler.InletExpansionAngle(hRel)
}

func InletPSAngle(hRel float64, profiler Profiler) float64 {
	return InletBendAngle(hRel, profiler) - InletPSExpansionAngle(hRel, profiler)
}

func InletSSAngle(hRel float64, profiler Profiler) float64 {
	return InletBendAngle(hRel, profiler) + InletSSExpansionAngle(hRel, profiler)
}

func OutletPSAngle(hRel float64, profiler Profiler) float64 {
	return OutletBendAngle(hRel, profiler) - OutletPSExpansionAngle(hRel, profiler)
}

func OutletSSAngle(hRel float64, profiler Profiler) float64 {
	return OutletBendAngle(hRel, profiler) + OutletSSExpansionAngle(hRel, profiler)
}

type profiler struct {
	behavior ProfilingBehavior

	geomGen geometry.BladingGeometryGenerator

	meanInletTriangle  states.VelocityTriangle
	meanOutletTriangle states.VelocityTriangle

	inletVelocityLaw  radial.VelocityLaw
	outletVelocityLaw radial.VelocityLaw

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

func (profiler *profiler) BladeNumber() int {
	var baRel = profiler.geomGen.Elongation()
	var lRelOut = profiler.geomGen.LRelOut()

	return common.RoundInt(math.Pi * baRel / lRelOut * 1 / profiler.approxTRel)
}

func (profiler *profiler) inletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.inletVelocityLaw.InletTriangle(
		profiler.meanInletTriangle,
		hRel,
		geometry.LRelIn(profiler.geomGen),
	)
}

func (profiler *profiler) outletTriangle(hRel float64) states.VelocityTriangle {
	return profiler.outletVelocityLaw.OutletTriangle(
		profiler.meanOutletTriangle,
		hRel,
		profiler.geomGen.LRelOut(),
	)
}
