package util

import (
	"fmt"
	"io"
	"os"

	"github.com/juju/errors"
)

type ErrorHandler struct {
	Out io.Writer
	PrintStackTrace bool
}

// The default error handler
var E = ErrorHandler{
	Out: os.Stderr,
	PrintStackTrace: true,
}

// New creates a new error variable
func (e *ErrorHandler) New(format string, a ...interface{}) error {
	ret := errors.NewErr(format, a...)
	ret.SetLocation(1)
	return &ret
}

// Annotate increases context information to the error
func (e *ErrorHandler) Annotate(err error, a ...interface{}) error {
	ret := errors.Annotate(err, fmt.Sprint(a...))
	if err, ok := ret.(*errors.Err); ok {
		err.SetLocation(1)
		ret = err
	}
	return ret
}

// Print writes the error message to predefined io.Writer
func (e *ErrorHandler) Print(err error, a ...interface{}) {
	fmt.Fprint(e.Out, "Error: ")
	if len (a) > 0 {
		fmt.Fprint(e.Out, a...)
		fmt.Fprint(e.Out, ": ")
	}
	fmt.Fprintf(e.Out, "%s\n", err)

	if e.PrintStackTrace {
		fmt.Fprintln(e.Out, "Error stack:")
		fmt.Fprintln(e.Out, errors.ErrorStack(err))
	}
}

func (e *ErrorHandler) Panic(err error, a ...interface{}) {
	e.Print(err, a...)
	panic("Irrecoverable error")
}
