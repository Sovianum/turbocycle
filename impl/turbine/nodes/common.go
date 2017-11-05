package nodes

import "github.com/Sovianum/turbocycle/core"

const (
	MassRateInput   = "MassRateInput"
	MassRateOutput  = "MassRateOutput"
	DimensionInput  = "DimensionInput"
	DimensionOutput = "DimensionOutput"
	VelocityInput   = "VelocityInput"
	VelocityOutput  = "VelocityOutput"
)

type MassRateChannel interface {
	MassRateSource
	MassRateSink
}

type MassRateSink interface {
	MassRateInput() core.Port
}

type MassRateSource interface {
	MassRateOutput() core.Port
}

type VelocityChannel interface {
	VelocitySink
	VelocitySource
}

type VelocitySink interface {
	VelocityInput() core.Port
}

type VelocitySource interface {
	VelocityOutput() core.Port
}

type DimensionChannel interface {
	DimensionSource
	DimensionSink
}

type DimensionSource interface {
	DimensionInput() core.Port
}

type DimensionSink interface {
	DimensionSource() core.Port
}
