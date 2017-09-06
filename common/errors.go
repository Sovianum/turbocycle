package common

import (
	"fmt"
	"errors"
)

func GetTypeError(expected string, got interface{}) error {
	return errors.New(fmt.Sprintf("Expected type %s, got %v", expected, got))
}
