package errors

import "testing"

// ensure that Wrap and Wrapf always return nil on a nil error or non-erroring code
// could get really weird for the caller.
func TestNilReturns(t *testing.T) {
	if err := Wrap(nil, "foo"); err != nil {
		t.Fatal("Wrap(nil) should have returned nil. got: ", err.Error())
	}

	if err := Wrapf(nil, "foo"); err != nil {
		t.Fatal("Wrapf(nil) should have returned nil. got: ", err.Error())
	}
}
