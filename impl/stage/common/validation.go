package common

import (
	"fmt"
	"math"
)

func NotNanValidator(x float64) error {
	if math.IsNaN(x) {
		return fmt.Errorf("nan obtained")
	}
	return nil
}
