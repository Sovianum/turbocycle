package constructive

import (
	"github.com/Sovianum/turbocycle/core"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/helpers/gdf"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"math"
)

type TurbineNode interface {
	core.Node
	nodes.ComplexGasSink
	nodes.PowerSource
	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut
	PiTStag() float64
	InputGas() gases.Gas
	LambdaOut() float64
}

type GasBiPole interface {
	core.Node
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
	var cp = gases.CpMean(node.InputGas(), node.TStagIn(), node.TStagOut(), nodes.DefaultN)
	var k = gases.KMean(node.InputGas(), node.TStagIn(), node.TStagOut(), nodes.DefaultN)
	var pi = gdf.Pi(node.LambdaOut(), k)
	var x = math.Pow(node.PiTStag()/pi, (1-k)/k)

	return cp * node.TStagIn() * (1 - x)
}
