package variator

type Variator interface {
	GetValue() float64
	SetValue(value float64)
}

func Add(variator Variator, delta float64) {
	variator.SetValue(variator.GetValue() + delta)
}

func FromPointer(ptr *float64) Variator {
	return FromCallables(
		func() float64 {
			return *ptr
		},
		func(f float64) {
			*ptr = f
		},
	)
}

func FromCallables(getter func() float64, setter func(float64)) Variator {
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
