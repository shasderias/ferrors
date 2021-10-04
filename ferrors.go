package ferrors

import (
	"fmt"
)

// Wrap returns the given error transparently wrapped, without additional context.
// As a special case, a nil argument results in a nil return.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return &wrapError{"", err, Caller(1)}
}

type appError struct {
	statusCode int
	msg        string
	err        error
	frame      Frame
}

func AppError(statusCode int, err error, msg string) error {
	if msg == "" && err == nil {
		return nil
	}

	return &appError{
		statusCode: statusCode,
		msg:        msg,
		err:        err,
		frame:      Caller(1),
	}
}

func AppErrorf(statusCode int, err error, format string, args ...interface{}) error {
	if format == "" && err == nil {
		return nil
	}

	return &appError{
		statusCode: statusCode,
		msg:        fmt.Sprintf(format, args...),
		err:        err,
		frame:      Caller(1),
	}
}

func (e appError) StatusCode() int {
	return e.statusCode
}

func (e appError) Msg() string {
	return e.msg
}

func (e appError) Error() string {
	switch {
	case e.msg == "" && e.err == nil:
		return ""
	case e.msg == "" && e.err != nil:
		return e.err.Error()
	case e.msg != "" && e.err == nil:
		return e.msg
	default:
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
}

func (e *appError) Format(s fmt.State, v rune) {
	// PATCH(robfig): Delegate to Error() unless + is specified.
	if s.Flag('+') {
		FormatError(e, s, v)
	} else {
		_, _ = s.Write([]byte(e.Error()))
	}
}

func (e *appError) FormatError(p Printer) (next error) {
	p.Print(e.msg)
	e.frame.Format(p)
	return e.err
}

func (e appError) Unwrap() error {
	return e.err
}

type ErrWithMsg interface {
	Msg() string
}

type ErrWithStatusCode interface {
	StatusCode() int
}
