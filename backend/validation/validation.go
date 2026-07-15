package validation

type Rule func() error

type Validator interface {
	Validate(any) error
}

type Errors struct {
	Items []Error
}

type Error struct {
	Field   string
	Message string
}

func (e *Errors) Add(field, msg string) {
	e.Items = append(e.Items, Error{Field: field, Message: msg})
}

func (e *Errors) HasErrors() bool {
	return len(e.Items) > 0
}

func (e *Errors) Error() string {
	if len(e.Items) == 0 {
		return ""
	}
	return e.Items[0].Message
}
