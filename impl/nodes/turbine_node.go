package nodes

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/gdf"
	"math"
)

type TurbineNode interface {
	core.Node
	TStagIn() float64
	TStagOut() float64
	PStagIn() float64
	PStagOut() float64
	Pit() float64
	InputGas() gases.Gas
	LambdaOut() float64
}

func TOut(node TurbineNode) float64 {
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)

	return tStag * gdf.Tau(node.LambdaOut(), k)
}

func POut(node TurbineNode) float64 {
	var pStag = node.PStagOut()
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)

	return pStag * gdf.Pi(node.LambdaOut(), k)
}

func DensityOut(node TurbineNode) float64 {
	return gases.Density(
		node.InputGas(), TOut(node), POut(node),
	)
}

func VelocityOut(node TurbineNode) float64 {
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)
	var r = node.InputGas().R()

	return gdf.ACrit(k, r, tStag) * node.LambdaOut()
}

func Ht(node TurbineNode) float64 {
	var cp = gases.CpMean(node.InputGas(), node.TStagIn(), node.TStagOut(), defaultN)
	var k = gases.KMean(node.InputGas(), node.TStagIn(), node.TStagOut(), defaultN)
	var x = math.Pow(node.Pit(), (1-k)/k)

	return cp * node.TStagIn() * (1 - x)
}
