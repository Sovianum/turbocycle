package profilers

import (
	"math"

	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

type GasProfiler interface {
	Density(hRel float64) float64
	Temperature(hRel float64) float64
	Pressure(hRel float64) float64
	CLambda(hRel float64) float64
	CMachNumber(hRel float64) float64
	CSpeedSound(hRel float64) float64
	KMean() float64
	ReactivityMean() float64
}

func Reactivity(hRel, hRelMean float64, inletProfiler, outletProfiler GasProfiler) float64 {
	var pi = PressureDrop(hRel, inletProfiler, outletProfiler)
	var piMean = PressureDrop(hRelMean, inletProfiler, outletProfiler)
	var k = inletProfiler.KMean()
	var reactivityRelation = (math.Pow(pi, (k-1)/k) - 1) / (math.Pow(piMean, (k-1)/k) - 1)
	var result = inletProfiler.ReactivityMean() * reactivityRelation
	return result
}

func PressureDrop(hRel float64, inletProfiler, outletProfiler GasProfiler) float64 {
	var inletPressure = inletProfiler.Pressure(hRel)
	var outletPressure = outletProfiler.Pressure(hRel)
	return inletPressure / outletPressure
}

func OutletGasProfiler(
	gas gases.Gas, tStag, pStag, reactivity float64, triangle states.VelocityTriangle, profiler Profiler,
) GasProfiler {
	return &gasProfiler{
		side:           NewOutletSide(profiler),
		hRelMean:       0.5,
		gas:            gas,
		tMean:          tStag,
		pMean:          pStag,
		reactivityMean: reactivity,
		triangleMean:   triangle,
	}
}

func InletGasProfiler(
	gas gases.Gas, tStag, pStag, reactivity float64, triangle states.VelocityTriangle, profiler Profiler,
) GasProfiler {
	return &gasProfiler{
		side:           NewInletSide(profiler),
		hRelMean:       0.5,
		gas:            gas,
		tMean:          tStag,
		pMean:          pStag,
		reactivityMean: reactivity,
		triangleMean:   triangle,
	}
}

type gasProfiler struct {
	side           Side
	hRelMean       float64
	gas            gases.Gas
	tMean          float64
	pMean          float64
	reactivityMean float64
	triangleMean   states.VelocityTriangle
}

func (gp *gasProfiler) ReactivityMean() float64 {
	return gp.reactivityMean
}

func (gp *gasProfiler) KMean() float64 {
	return gases.K(gp.gas, gp.tMean)
}

func (gp *gasProfiler) Density(hRel float64) float64 {
	var result = gp.Pressure(hRel) / (gp.gas.R() * gp.Temperature(hRel))
	return result
}

func (gp *gasProfiler) Temperature(hRel float64) float64 {
	var lambda = gp.CLambda(hRel)
	var lambdaMean = gp.CLambda(gp.hRelMean)
	var k = gases.K(gp.gas, gp.tMean)
	var result = gp.pMean / gdf.Tau(lambdaMean, k) * gdf.Tau(lambda, k)
	return result
}

func (gp *gasProfiler) Pressure(hRel float64) float64 {
	var lambda = gp.CLambda(hRel)
	var lambdaMean = gp.CLambda(gp.hRelMean)
	var k = gases.K(gp.gas, gp.tMean)
	var piMean = gdf.Pi(lambdaMean, k)
	var pi = gdf.Pi(lambda, k)

	var result = gp.pMean / piMean * pi
	return result
}

func (gp *gasProfiler) CLambda(hRel float64) float64 {
	var k = gases.K(gp.gas, gp.tMean)
	var mach = gp.CMachNumber(hRel)
	var result = gdf.Lambda(mach, k)
	return result
}

func (gp *gasProfiler) CMachNumber(hRel float64) float64 {
	var triangle = gp.side.Triangle(hRel)
	var speedSound = gp.CSpeedSound(hRel)
	var result = triangle.C() / speedSound
	return result
}

func (gp *gasProfiler) CSpeedSound(hRel float64) float64 {
	var c = gp.side.Triangle(hRel).C()
	var cMean = gp.triangleMean.C()

	var t = gp.tMean - (c*c-cMean*cMean)/(2*gp.gas.Cp(gp.tMean))
	var result = math.Sqrt(gp.gas.R() * gases.K(gp.gas, gp.tMean) * t)
	return result
}
