package errors

type Kind int

const (
	KindInternal    Kind = iota
	KindValidation
	KindNotFound
	KindUnauthorized
	KindForbidden
	KindConflict
	KindTimeout
	KindExternal
)

type Error struct {
	Kind    Kind
	Message string
	Op      string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, op, msg string) *Error {
	return &Error{Kind: kind, Op: op, Message: msg}
}

func Wrap(kind Kind, op, msg string, err error) *Error {
	return &Error{Kind: kind, Op: op, Message: msg, Err: err}
}
