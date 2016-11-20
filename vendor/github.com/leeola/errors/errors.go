package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

func New(s string) error {
	return &errWrap{
		Msg:      s,
		SumStack: []string{callerLine() + ": " + s},
	}
}

func Cause(err error) error {
	cErr, ok := err.(causer)
	if !ok {
		return err
	}

	if cause := cErr.Cause(); cause != nil {
		return cause
	}

	return err
}

func Errorf(f string, s ...interface{}) error {
	msg := fmt.Sprintf(f, s...)
	return &errWrap{
		Msg:      msg,
		SumStack: []string{callerLine() + ": " + msg},
	}
}

func Println(err error) {
	if err == nil {
		return
	}

	fmt.Print(Sprintln(err))
}

func Sprintln(err error) string {
	if err == nil {
		return ""
	}

	sErr, ok := err.(*errWrap)
	if !ok {
		return err.Error() + "\n"
	}

	return strings.Join(sErr.SumStack, "\n") + "\n"
}

func Wrap(err error, s string) error {
	if err == nil {
		return nil
	}

	return wrap(err, callerLine(), s)
}

func Wrapf(err error, f string, s ...interface{}) error {
	if err == nil {
		return nil
	}

	return wrap(err, callerLine(), fmt.Sprintf(f, s...))
}

func wrap(err error, caller string, s string) error {
	sErr, ok := err.(*errWrap)

	// construct the cascaded error line
	cascadeErrLine := s + ": " + err.Error()

	// If it's not a wrapped error, construct a new one
	if !ok {
		// the previous err is *not* a wrapped error, so include the caller err.Error()
		// message in the stackLine.
		stackLine := caller + ": " + cascadeErrLine

		return &errWrap{
			err:      err,
			Msg:      cascadeErrLine,
			SumStack: []string{stackLine},
		}
	}

	// the given error is a wrapped error, *do not* include the
	// err.Error(), as that will consist of the cascade message.
	stackLine := caller + ": " + s

	// the given error is a wrapped error, modify it to the latest error information
	sErr.Msg = cascadeErrLine
	sErr.SumStack = append(sErr.SumStack, stackLine)
	return sErr
}

// callerLine returns a short path and the line number of the caller
//
// Currently it's returning the two directories above the file for brevity. In
// the future we may want to display the full path relative to the $GOPATH.
func callerLine() string {
	// p comes out as the absolute path, it needs to be trimmed.
	_, p, l, _ := runtime.Caller(2)
	// TODO(leeola): there has to be a better way to write this...
	p, f := path.Dir(p), path.Base(p)
	p, parent := path.Dir(p), path.Base(p)
	p, great := path.Dir(p), path.Base(p)
	p = path.Join(great, parent, f)
	return fmt.Sprintf("%s:%d", p, l)
}

type causer interface {
	Cause() error
}

type errWrap struct {
	// The original error that this error is wrapping. Stored to retrieve the
	// original as needed.
	err error

	Msg      string
	SumStack []string
}

func (e *errWrap) Error() string {
	return e.Msg
}

func (e *errWrap) Errors() []string {
	return e.SumStack[:]
}

func (e *errWrap) Cause() error {
	if e.err == nil {
		return e.err
	}
	return e
}
