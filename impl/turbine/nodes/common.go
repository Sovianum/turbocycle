package nodes

import "github.com/Sovianum/turbocycle/core/graph"

type MassRateChannel interface {
	MassRateSource
	MassRateSink
}

type MassRateSink interface {
	MassRateInput() graph.Port
}

type MassRateSource interface {
	MassRateOutput() graph.Port
}

type VelocityChannel interface {
	VelocitySink
	VelocitySource
}

type VelocitySink interface {
	VelocityInput() graph.Port
}

type VelocitySource interface {
	VelocityOutput() graph.Port
}

type DimensionChannel interface {
	DimensionSource
	DimensionSink
}

type DimensionSource interface {
	DimensionInput() graph.Port
}

type DimensionSink interface {
	DimensionSource() graph.Port
}
