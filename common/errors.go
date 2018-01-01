package common

import (
	"fmt"
)

func GetTypeError(expected string, got interface{}) error {
	return fmt.Errorf("expected type %s, got %v", expected, got)
}
