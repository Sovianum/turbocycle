package nodes

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/Sovianum/turbocycle/gases"
	"fmt"
)

func TestCompressorNode_PStagIn(t *testing.T) {
	var compressor = NewCompressorNode(0.89, 6, 0.95)
	assert.NotNil(t, compressor)

	var inputState = states.NewGasPortState(gases.GetAir(), 500, 100, 1)
	compressor.GasInput().SetState(inputState)

	assert.Equal(t, compressor.PStagIn(), float64(100))
}

func TestCompressorNode_Process(t *testing.T) {
	var compressor = NewCompressorNode(0.89, 6, 0.95)
	assert.NotNil(t, compressor)

	var inputState = states.NewGasPortState(gases.GetAir(), 288, 100, 1)
	compressor.GasInput().SetState(inputState)

	compressor.Process()

	assert.Equal(t, compressor.PStagOut(), float64(600))

	fmt.Println(compressor.GasOutput().GetState())
}
