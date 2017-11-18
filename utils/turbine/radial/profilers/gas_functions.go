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
	return inletProfiler.ReactivityMean() * reactivityRelation
}

func PressureDrop(hRel float64, inletProfiler, outletProfiler GasProfiler) float64 {
	return inletProfiler.Pressure(hRel) / outletProfiler.Pressure(hRel)
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
	return gp.Pressure(hRel) / (gp.gas.R() * gp.Temperature(hRel))
}

func (gp *gasProfiler) Temperature(hRel float64) float64 {
	var lambda = gp.CLambda(hRel)
	var lambdaMean = gp.CLambda(gp.hRelMean)
	var k = gases.K(gp.gas, gp.tMean)

	return gp.pMean / gdf.Tau(lambdaMean, k) * gdf.Tau(lambda, k)
}

func (gp *gasProfiler) Pressure(hRel float64) float64 {
	var lambda = gp.CLambda(hRel)
	var lambdaMean = gp.CLambda(gp.hRelMean)
	var k = gases.K(gp.gas, gp.tMean)

	return gp.pMean / gdf.Pi(lambdaMean, k) * gdf.Pi(lambda, k)
}

func (gp *gasProfiler) CLambda(hRel float64) float64 {
	var k = gases.K(gp.gas, gp.tMean)
	return gdf.Lambda(
		gp.CMachNumber(hRel),
		k,
	)
}

func (gp *gasProfiler) CMachNumber(hRel float64) float64 {
	return gp.side.Triangle(hRel).C() / gp.CSpeedSound(hRel)
}

func (gp *gasProfiler) CSpeedSound(hRel float64) float64 {
	var c = gp.side.Triangle(hRel).C()
	var cMean = gp.triangleMean.C()

	var t = gp.tMean + (c*c/2+cMean*cMean/2)/(2*gp.gas.Cp(gp.tMean))
	return math.Sqrt(gp.gas.R() * gases.K(gp.gas, gp.tMean) * t)
}
