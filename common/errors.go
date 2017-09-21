package common

import (
	"fmt"
)

func GetTypeError(expected string, got interface{}) error {
	return fmt.Errorf("Expected type %s, got %v", expected, got)
}
