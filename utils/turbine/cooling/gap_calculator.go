package cooling

import (
	"fmt"
	"math"

	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
)

func NewGapCalculator(
	cooler, gas gases.Gas,
	ca, pGas float64,
	bladeGeometry geometry.BladingGeometry,
	profile profiles.BladeProfile,
	wallThk float64,
	lambdaM float64,
	nuGasFunc func(re float64) float64,

	tGas, tWallOuter, tCoolerInlet float64,
) GapCalculator {
	var height = geometry.Height(0, bladeGeometry)
	var perimeter = profiles.Perimeter(profile)
	var area = perimeter * height

	return &gapCalculator{
		cooler: cooler,
		gas:    gas,

		ca:   ca,
		pGas: pGas,

		ba:          geometry.ChordProjection(bladeGeometry),
		bladeLength: height,
		perimeter:   perimeter,
		bladeArea:   area,
		wallThk:     wallThk,

		lambdaM: lambdaM,

		nuGasFunc: nuGasFunc,

		tGas:         tGas,
		tWallOuter:   tWallOuter,
		tCoolerInlet: tCoolerInlet,
	}
}

type DataPack struct {
	Err error

	BladeLength     float64
	ChordProjection float64
	BladeArea       float64
	Perimeter       float64
	WallThk         float64

	TGas       float64
	CaGas      float64
	DensityGas float64
	MuGas      float64
	LambdaGas  float64

	MassRateCooler float64

	ReGas    float64
	NuGas float64
	NuCoef float64
	AlphaGas float64

	BladeHeat float64

	TAir0 float64

	TDrop      float64
	TMean      float64
	TWallInner float64
	TWallOuter float64

	LambdaM float64

	DComplex   float64
	EpsComplex float64

	AirGap float64
}

type GapCalculator interface {
	GetPack(coolerMassRate float64) DataPack
}

type gapCalculator struct {
	cooler gases.Gas
	gas    gases.Gas

	ca   float64
	pGas float64

	ba          float64
	bladeLength float64
	perimeter   float64
	bladeArea   float64
	wallThk     float64

	lambdaM float64

	nuGasFunc func(re float64) float64

	tGas         float64
	tWallOuter   float64
	tCoolerInlet float64
}

func (calc *gapCalculator) GetPack(coolerMassRate float64) DataPack {
	var pack = new(DataPack)
	pack.MassRateCooler = coolerMassRate
	calc.initPack(pack)

	calc.reGas(pack)
	calc.alphaGas(pack)
	calc.bladeHeat(pack)
	calc.tDrop(pack)
	calc.tMean(pack)
	calc.tWallInner(pack)
	calc.epsComplex(pack)
	calc.dComplex(coolerMassRate, pack)
	calc.airGap(coolerMassRate, pack)
	return *pack
}

func (calc *gapCalculator) airGap(coolerMassRate float64, pack *DataPack) {
	var airGap = pack.EpsComplex * math.Pow(coolerMassRate, 0.8) * (pack.DComplex - calc.bladeArea/(2*calc.cooler.Cp(calc.tCoolerInlet)*coolerMassRate))
	if airGap < 0 {
		pack.Err = fmt.Errorf("D complex is less than term3 (mass rate = %f)", coolerMassRate)
		return
	}
	var fComplex = calc.bladeArea / (2 * calc.cooler.Cp(calc.tCoolerInlet) * coolerMassRate)
	pack.AirGap = pack.EpsComplex * math.Pow(coolerMassRate, 0.8) * (pack.DComplex - fComplex)
}

func (calc *gapCalculator) dComplex(coolerMassRate float64, pack *DataPack) {
	var term1 = 1 / pack.AlphaGas * ((calc.tGas-calc.tCoolerInlet)/(calc.tGas-calc.tWallOuter) - 1)
	var term2 = -calc.wallThk / calc.lambdaM

	var dComplex = term1 + term2
	pack.DComplex = dComplex
}

func (calc *gapCalculator) epsComplex(pack *DataPack) {
	var mu = calc.cooler.Mu(calc.tCoolerInlet)
	var lambda = calc.cooler.Lambda(calc.tCoolerInlet)
	pack.EpsComplex = 0.01 * lambda * math.Pow(1/(calc.bladeLength*mu), 0.8) // todo maybe need to abstract out
}

func (calc *gapCalculator) tWallInner(pack *DataPack) {
	pack.TWallInner = calc.tWallOuter - pack.TDrop
}

func (calc *gapCalculator) tMean(pack *DataPack) {
	pack.TMean = calc.tWallOuter - pack.TDrop/2
}

func (calc *gapCalculator) tDrop(pack *DataPack) {
	pack.TDrop = pack.BladeHeat * calc.wallThk / (calc.bladeArea * calc.lambdaM)
}

func (calc *gapCalculator) bladeHeat(pack *DataPack) {
	pack.BladeHeat = pack.AlphaGas * calc.perimeter * calc.bladeLength * (calc.tGas - calc.tWallOuter)
}

func (calc *gapCalculator) alphaGas(pack *DataPack) {
	var lambda = calc.gas.Lambda(calc.tGas)
	pack.LambdaGas = lambda
	var nu = calc.nuGasFunc(pack.ReGas)
	pack.NuGas = nu
	pack.AlphaGas = lambda / calc.ba * nu
}

func (calc *gapCalculator) reGas(pack *DataPack) {
	var density = gases.Density(calc.gas, calc.tGas, calc.pGas)
	pack.DensityGas = density
	var mu = calc.gas.Mu(calc.tGas)
	pack.MuGas = mu
	pack.ReGas = calc.ca * calc.ba * density / mu
}

func (calc *gapCalculator) initPack(pack *DataPack) {
	pack.BladeLength = calc.bladeLength
	pack.ChordProjection = calc.ba
	pack.BladeArea = calc.bladeArea
	pack.Perimeter = calc.perimeter
	pack.WallThk = calc.wallThk

	pack.TGas = calc.tGas
	pack.CaGas = calc.ca
	pack.TAir0 = calc.tCoolerInlet
	pack.LambdaM = calc.lambdaM
	pack.TWallOuter = calc.tWallOuter
}
