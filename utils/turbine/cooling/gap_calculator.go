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
	err error

	re       float64
	alphaGas float64

	bladeHeat float64

	tDrop      float64
	tMean      float64
	tWallInner float64

	dComplex   float64
	epsComplex float64

	airGap float64
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
	calc.re(pack)
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
	if pack.err != nil {
		pack.err = fmt.Errorf("%s: airGap", pack.err.Error())
		return
	}
	pack.airGap = pack.epsComplex * math.Pow(coolerMassRate, 0.8) * (pack.dComplex - calc.bladeArea/(2*calc.cooler.Cp(calc.tGas)*coolerMassRate))
}

func (calc *gapCalculator) dComplex(coolerMassRate float64, pack *DataPack) {
	var term1 = 1 / pack.alphaGas * ((calc.tGas-calc.tCoolerInlet)/(calc.tGas-calc.tWallOuter) - 1)
	var term2 = -calc.wallThk / calc.lambdaM
	var term3 = -calc.bladeArea / (2 * calc.cooler.Cp(calc.tGas) * coolerMassRate)

	var dComplex = term1 + term2 + term3
	if dComplex <= math.Abs(term3) {
		pack.err = fmt.Errorf("D complex is less than term3")
		return
	}
	pack.dComplex = dComplex
}

func (calc *gapCalculator) epsComplex(pack *DataPack) {
	var mu = calc.cooler.Mu(calc.tCoolerInlet)
	var lambda = calc.cooler.Lambda(calc.tCoolerInlet)
	pack.epsComplex = 0.01 * lambda * math.Pow(1/(calc.bladeLength*mu), 0.8) // todo maybe need to abstract out
}

func (calc *gapCalculator) tWallInner(pack *DataPack) {
	pack.tWallInner = calc.tWallOuter - pack.tDrop
}

func (calc *gapCalculator) tMean(pack *DataPack) {
	pack.tMean = calc.tWallOuter - pack.tDrop/2
}

func (calc *gapCalculator) tDrop(pack *DataPack) {
	pack.tDrop = pack.bladeHeat * calc.wallThk / (calc.bladeArea * calc.lambdaM)
}

func (calc *gapCalculator) bladeHeat(pack *DataPack) {
	pack.bladeHeat = pack.alphaGas * calc.perimeter * calc.bladeLength * (calc.tGas - calc.tWallOuter)
}

func (calc *gapCalculator) alphaGas(pack *DataPack) {
	pack.alphaGas = calc.gas.Lambda(calc.tGas) / calc.ba * calc.nuGasFunc(pack.re)
}

func (calc *gapCalculator) re(pack *DataPack) {
	pack.re = calc.ca * calc.ba * gases.Density(calc.gas, calc.tGas, calc.pGas) / calc.gas.Mu(calc.tGas)
}
