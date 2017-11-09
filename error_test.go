package util

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

func structEquals(a, b interface{}) bool {
	return spew.Sdump(a) == spew.Sdump(b)
}

func diffStr(a, b interface{}) (ret string) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(spew.Sdump(a)),
		B:        difflib.SplitLines(spew.Sdump(b)),
		FromFile: "Expected",
		ToFile:   "Received",
		Context:  3,
	}

	ret, _ = difflib.GetUnifiedDiffString(diff)
	return
}

type ehOp interface {
	run(*ErrorHandler, error) error
}

type ehFunc func(*ErrorHandler, error) error

func (f ehFunc) run(eh *ErrorHandler, err error) error {
	return f(eh, err)
}

func TestErrorHandler(t *testing.T) {
	n := func(msg string) ehFunc {
		return func(e *ErrorHandler, err error) error {
			return e.New(msg)
		}
	}

	a := func(msg string) ehFunc {
		return func(e *ErrorHandler, err error) error {
			return e.Annotate(err, msg)
		}
	}

	type args struct {
	}
	tests := []struct {
		name            string
		ops             []ehOp
		out             string
		PrintStackTrace bool
	}{
		{"Simple error", []ehOp{n("first")}, "first", false},
		{"Annotating error", []ehOp{n("first"), a("sec")}, "sec: first", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			var err error = nil
			e := &ErrorHandler{
				Out:             buf,
				PrintStackTrace: tt.PrintStackTrace,
			}

			for _, op := range tt.ops {
				err = op.run(e, err)
			}

			if !structEquals(tt.out, err.Error()) {
				t.Errorf("Expected error message differs:\n %s", diffStr(tt.out, err.Error()))

			}

			e.Print(err, "a")
			out := "Error: a: " + tt.out + "\n"

			if !structEquals(out, buf.String()) {
				t.Errorf("Expected printed error message differs:\n %s", diffStr(out, buf.String()))
			}
		})
	}
}

type testErr struct {
	msg string
}

func (t *testErr) Error() string {
	return t.msg
}

type testOp interface {
	run(*ErrorList)
}

type testFunc func(*ErrorList)

func (t testFunc) run(d *ErrorList) {
	t(d)
}

func TestErrorList(t *testing.T) {
	ae := func(errmsg string) testFunc {
		return func(e *ErrorList) {
			e.Append(&testErr{errmsg})
		}
	}
	tests := []struct {
		name    string
		message string
		ops     []testOp
		result  string
		isEmpty bool
	}{
		{"Empty list", "empty", []testOp{}, "", true},
		{"One error", "one", []testOp{ae("a")}, "Error: one: Error 1: a; ", false},
		{"Two errors", "two", []testOp{ae("a"), ae("b")}, "Error: two: Error 1: a; Error 2: b; ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			el := NewErrorList(tt.message)

			for _, op := range tt.ops {
				op.run(el)
			}

			gotRet := el.Error()
			if !structEquals(gotRet, tt.result) {
				t.Errorf("Expected error message differs:\n %s", diffStr(gotRet, tt.result))
			}

			if el.IsEmpty() != tt.isEmpty {
				t.Errorf("Expected to be empty: \"%v\" Reported empty: \"%v\"", tt.isEmpty, el.IsEmpty())
			}
		})
	}
}
