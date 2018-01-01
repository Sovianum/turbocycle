package schemes

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
)

type Scheme interface {
	GetNetwork() (graph.Network, graph.GraphError)
	GetSpecificPower() float64
	GetFuelMassRateRel() float64
	GetQLower() float64
}

type SingleCompressor interface {
	Compressor() constructive.CompressorNode
}

type DoubleCompressor interface {
	LowPressureCompressor() constructive.CompressorNode
	HighPressureCompressor() constructive.CompressorNode
}

func GetMassRate(power float64, scheme Scheme) float64 {
	return power / scheme.GetSpecificPower()
}

func GetSpecificFuelRate(scheme Scheme) float64 {
	return 3600 * scheme.GetFuelMassRateRel() / scheme.GetSpecificPower()
}

func GetEfficiency(scheme Scheme) float64 {
	return scheme.GetSpecificPower() / (scheme.GetFuelMassRateRel() * scheme.GetQLower())
}
