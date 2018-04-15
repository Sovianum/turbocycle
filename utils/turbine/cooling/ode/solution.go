package ode

type Solution interface {
	Build() (xArr, yArr []float64)
	Step() float64
}
