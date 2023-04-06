package errors

type IError interface {
	Error() string
}

type Error struct {
	string
}

func (err Error) Error() string { return err.string }

var (
	ErrorContextCanceled = &Error{"context canceled"}

	ErrorData = &Error{"data error"}
)

func NewError(err error) *Error {
	return &Error{err.Error()}
}
