package errors

import (
	"github.com/rotisserie/eris"
)

var (
	ErrorContextCanceled = eris.New("context canceled")
	ErrorData            = eris.New("data error")
)
