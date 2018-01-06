package geom

type Transformable interface {
	Transform(t Transformation)
}
