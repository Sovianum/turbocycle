package common

type DiscreteFunc func(id int) float64

func (f DiscreteFunc) SetValue(x int, y float64) DiscreteFunc {
	return func(id int) float64 {
		if id == x {
			return y
		}
		return f(id)
	}
}

func FromConst(c float64) DiscreteFunc {
	return func(id int) float64 {
		return c
	}
}

func FromArr(arr []float64, extend bool) DiscreteFunc {
	return func(id int) float64 {
		localID := id
		if id >= len(arr) && extend {
			localID = len(arr) - 1
		}
		return arr[localID]
	}
}

func FromDistribution(d Func1D) DiscreteFunc {
	return func(id int) float64 {
		return d(float64(id))
	}
}
