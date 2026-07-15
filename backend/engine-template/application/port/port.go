package port

type InputPort interface {
	Execute(any) (any, error)
}

type OutputPort interface {
	Present(any) any
}
