package math

type Variator interface {
	GetValue() float64
	SetValue(value float64)
}

func Add(variator variator, delta float64) {
	variator.SetValue(variator.GetValue() + delta)
}

func VariatorFromPointer(ptr *float64) Variator {
	return VariatorFromCallables(
		func() float64 {
			return *ptr
		},
		func(f float64) {
			*ptr = f
		},
	)
}

func VariatorFromCallables(getter func() float64, setter func(float64)) Variator {
	return &variator{
		getter: getter,
		setter: setter,
	}
}

type variator struct {
	getter func() float64
	setter func(float64)
}

func (v *variator) GetValue() float64 {
	return v.getter()
}

func (v *variator) SetValue(value float64) {
	v.setter(value)
}
