package envy

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPatch(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "returns default object and error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"
				ObjectChannels[functionName] = make(chan interface{}, 100)
				ErrorChannels[functionName] = make(chan error, 100)
				DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				DefaultErrors[functionName] = fmt.Errorf("ruh roh")

				obj := GetObject(functionName)
				g.Expect(obj).To(BeAssignableToTypeOf(struct{ Foo string }{}))
				g.Expect(obj.(struct{ Foo string })).To(Equal(struct{ Foo string }{Foo: "bar"}))

				err := GetError(functionName)
				g.Expect(err).To(MatchError("ruh roh"))
			},
		},
		{
			testCase: "returns pushed object and error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"

				ObjectChannels[functionName] = make(chan interface{}, 100)
				ErrorChannels[functionName] = make(chan error, 100)
				DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				DefaultErrors[functionName] = nil

				err := AddObjectReturns(functionName, struct{ Foo string }{Foo: "foo"})
				g.Expect(err).NotTo(HaveOccurred())

				err = AddErrorReturns(functionName, fmt.Errorf("some error"))
				g.Expect(err).NotTo(HaveOccurred())

				obj := GetObject(functionName)
				g.Expect(obj).To(BeAssignableToTypeOf(struct{ Foo string }{}))
				g.Expect(obj.(struct{ Foo string })).To(Equal(struct{ Foo string }{Foo: "foo"}))

				err = GetError(functionName)
				g.Expect(err).To(MatchError("some error"))
			},
		},
		{
			testCase: "error for function not found",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"

				ObjectChannels[functionName] = make(chan interface{}, 100)
				ErrorChannels[functionName] = make(chan error, 100)
				DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				DefaultErrors[functionName] = nil

				err := AddObjectReturns("notafunction", struct{ Foo string }{Foo: "foo"})
				g.Expect(err).To(HaveOccurred())

				err = AddErrorReturns("notafunction", fmt.Errorf("some error"))
				g.Expect(err).To(HaveOccurred())

			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
