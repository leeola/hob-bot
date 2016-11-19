package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

func New(s string) error {
	return &sumErr{
		msg:      s,
		SumStack: []string{callerLine() + ": " + s},
	}
}

func Errorf(f string, s ...interface{}) error {
	msg := fmt.Sprintf(f, s...)
	return &sumErr{
		msg:      msg,
		SumStack: []string{callerLine() + ": " + msg},
	}
}

func Sum(err error, s string) error {
	return sum(err, callerLine(), s)
}

func Sumf(err error, f string, s ...interface{}) error {
	return sum(err, callerLine(), fmt.Sprintf(f, s...))
}

func sum(err error, caller string, s string) error {
	if err == nil {
		return nil
	}

	sErr, ok := err.(*sumErr)

	// construct the cascaded error line
	cascadeErrLine := err.Error() + ": " + s

	// If it's not a sumerror, return a new one
	if !ok {
		// the previous err is *not* a sumErr, so include the caller err.Error() message
		// in the stackLine.
		stackLine := caller + ": " + cascadeErrLine

		return &sumErr{
			msg:      cascadeErrLine,
			SumStack: []string{stackLine},
		}
	}

	// the previous error (sum'd/wrap'd error) is a sumErr, *do not* include the
	// err.Error(), as that will consist of the cascade message.
	stackLine := caller + ": " + s

	return &sumErr{
		msg:      cascadeErrLine,
		SumStack: append(sErr.SumStack, stackLine),
	}
}

// TODO(leeola): Write a simple wrap method for embedding errors
func Wrap(err error, s string) error {
	return sum(err, callerLine(), s)
}

// TODO(leeola): Write a simple wrap method for embedding errors
func Wrapf(err error, f string, s ...interface{}) error {
	return sum(err, callerLine(), fmt.Sprintf(f, s...))
}

func callerLine() string {
	_, p, l, _ := runtime.Caller(2)

	p, f := path.Dir(p), path.Base(p)
	p, parent := path.Dir(p), path.Base(p)
	p, great := path.Dir(p), path.Base(p)
	p = path.Join(great, parent, f)

	// TODO(leeola): trim to project path
	return fmt.Sprintf("%s:%d", p, l)
}

func Sprintln(err error) string {
	sErr, ok := err.(*sumErr)
	if !ok {
		return err.Error() + "\n"
	}

	return strings.Join(sErr.SumStack, "\n") + "\n"
}

func Println(err error) {
	fmt.Print(Sprintln(err))
}

type sumErr struct {
	// The original error if this error is wrapping. Stored to retrieve the
	// original as needed.
	err error

	msg      string
	SumStack []string
}

func (e *sumErr) Error() string {
	return e.msg
}

func (e *sumErr) Errors() []string {
	return e.SumStack[:]
}
