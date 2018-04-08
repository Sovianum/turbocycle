package geometry

import (
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
)

func TestAngles(t *testing.T) {
	width := 107.
	gapFactor := 0.1
	d0Out := 897.
	d0In := 494.
	innerAngle := common.ToRadians(25)
	geom := NewBladingGeometry(
		width*(1-gapFactor),
		width*gapFactor,
		NewAxialProfileLine(0, d0In, innerAngle),
		NewAxialProfileLine(0, d0Out, 0), 0.5,
	)

	assert.InDelta(t, 296.894*2, geom.InnerProfile().Diameter(width), 1e-1)
}
