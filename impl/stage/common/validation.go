package common

import (
	"fmt"
	"math"
)

type Validator func(x float64) error

var ComplexPositiveValidator = Validator(PositiveValidator).Combine(FiniteValidator)

var FiniteValidator = Validator(NotInfValidator).Combine(NotNanValidator)

func (v Validator) Combine(another Validator) Validator {
	return func(x float64) error {
		if e1 := v(x); e1 != nil {
			return e1
		}
		if e2 := another(x); e2 != nil {
			return e2
		}
		return nil
	}
}

func PositiveValidator(x float64) error {
	if x <= 0 {
		return fmt.Errorf("obtained non-positive value")
	}
	return nil
}

func NotInfValidator(x float64) error {
	if math.IsInf(x, 0) {
		return fmt.Errorf("inf obtained")
	}
	return nil
}

func NotNanValidator(x float64) error {
	if math.IsNaN(x) {
		return fmt.Errorf("nan obtained")
	}
	return nil
}
