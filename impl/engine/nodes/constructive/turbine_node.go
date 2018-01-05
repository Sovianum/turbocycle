package constructive

import (
	"math"

	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/material/gases"
)

type TurbineNode interface {
	graph.Node
	nodes.PowerSource
	nodes.PressureIn
	nodes.PressureOut
	nodes.TemperatureIn
	nodes.TemperatureOut

	nodes.GasChannel
	nodes.PressureChannel
	nodes.TemperatureChannel
	nodes.MassRateSource

	PiTStag() float64
	InputGas() gases.Gas
	Eta() float64
	LSpecific() float64
	MassRateRel() float64
	LeakMassRateRel() float64
	CoolMassRateRel() float64
}

type StaticTurbineNode interface {
	TurbineNode
	LambdaOut() float64
}

func TurbineLabour(node TurbineNode) float64 {
	return node.MassRateRel() * node.LSpecific()
}

type GasBiPole interface {
	graph.Node
}

func TOut(node StaticTurbineNode) float64 {
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)

	return tStag * gdf.Tau(node.LambdaOut(), k)
}

func POut(node StaticTurbineNode) float64 {
	var pStag = node.PStagOut()
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)

	return pStag * gdf.Pi(node.LambdaOut(), k)
}

func DensityOut(node StaticTurbineNode) float64 {
	return gases.Density(
		node.InputGas(), TOut(node), POut(node),
	)
}

func VelocityOut(node StaticTurbineNode) float64 {
	var tStag = node.TStagOut()
	var k = gases.K(node.InputGas(), tStag)
	var r = node.InputGas().R()

	return gdf.ACrit(k, r, tStag) * node.LambdaOut()
}

func Ht(node StaticTurbineNode) float64 {
	var cp = gases.CpMean(node.InputGas(), node.TStagIn(), node.TStagOut(), nodes.DefaultN)
	var k = gases.KMean(node.InputGas(), node.TStagIn(), node.TStagOut(), nodes.DefaultN)
	var pi = gdf.Pi(node.LambdaOut(), k)
	var x = math.Pow(node.PiTStag()/pi, (1-k)/k)

	return cp * node.TStagIn() * (1 - x)
}
