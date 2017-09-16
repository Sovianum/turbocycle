package constructive

import (
	"github.com/Sovianum/turbocycle/impl/states"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransmissionNode(t *testing.T) {
	var node = NewTransmissionNode(0.99)
	var inputState = states.NewPowerPortState(100)

	node.PowerInput().SetState(inputState)
	var err = node.Process()
	assert.Nil(t, err)

	var outputState = node.PowerOutput().GetState().(states.PowerPortState)
	assert.Equal(t, 100/0.99, outputState.LSpecific)
}
