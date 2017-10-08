package schemes

import "github.com/Sovianum/turbocycle/core"

type Scheme interface {
	GetNetwork() core.Network
}
