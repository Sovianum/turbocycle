package nodes

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/gases"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
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

	var expectedTemperature = tA * (1 + 1/etaC*(math.Pow(piC, (kAir-1)/kAir)-1))
	assert.True(
		t,
		common.ApproxEqual(expectedTemperature, compressor.TStagOut(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedTemperature, compressor.TStagOut()),
	)

	var cp = gases.CpMean(gases.GetAir(), compressor.TStagIn(), compressor.TStagOut(), defaultN)
	var expectedLabour = cp * (compressor.TStagOut() - compressor.TStagIn())
	assert.True(
		t,
		common.ApproxEqual(expectedLabour, compressor.LSpecific(), 0.01),
		fmt.Sprintf("Expected %f, got %f", expectedLabour, compressor.LSpecific()),
	)
}

func getTestCompressor() CompressorNode {
	var compressor = NewCompressorNode(etaC, piC, 0.05)
	var gasState = states.NewGasPortState(gases.GetAir(), tA, pA, 1)
	compressor.GasInput().SetState(gasState)
	return compressor
}
