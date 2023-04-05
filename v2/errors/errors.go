package errors

type Error struct{ string }

func (err Error) Error() string { return err.string }

var (
	ErrorContextCanceled = &Error{"context canceled"}
)

func NewError(err error) *Error {
	return &Error{err.Error()}
}
