package constructive

import (
	"fmt"
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
)

const (
	piC  = 6
	etaC = 0.86
	pA   = 1e5
	tA   = 288
	kAir = 1.4
)

func TestCompressorNode_Process(t *testing.T) {
	var compressor = getTestCompressor()
	assert.NotNil(t, compressor)
	compressor.Process()

	var expectedPressure = float64(piC * pA)
	assert.Equal(t, compressor.PStagOut(), expectedPressure)

	var etaAd = (math.Pow(piC, (kAir-1)/kAir) - 1) / (math.Pow(piC, (kAir-1)/(kAir*etaC)) - 1)
	var expectedTemperature = tA * (1 + 1/etaAd*(math.Pow(piC, (kAir-1)/kAir)-1))
	assert.True(
		t,
		common.ApproxEqual(expectedTemperature, compressor.TStagOut(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedTemperature, compressor.TStagOut()),
	)

	var cp = gases.CpMean(gases.GetAir(), compressor.TStagIn(), compressor.TStagOut(), nodes.DefaultN)
	var expectedLabour = cp * (compressor.TStagOut() - compressor.TStagIn())
	assert.True(
		t,
		common.ApproxEqual(expectedLabour, compressor.LSpecific(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedLabour, compressor.LSpecific()),
	)
}

func getTestCompressor() CompressorNode {
	var compressor = NewCompressorNode(etaC, piC, 0.05)
	var gasState = states.NewComplexGasPortState(gases.GetAir(), tA, pA, 1)
	compressor.ComplexGasInput().SetState(gasState)
	return compressor
}
