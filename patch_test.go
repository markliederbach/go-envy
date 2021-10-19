package envy_test

import (
	"fmt"
	"testing"

	"github.com/markliederbach/go-envy"
	. "github.com/onsi/gomega"
)

func TestPatch(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "returns default object and error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"
				envy.ObjectChannels[functionName] = make(chan interface{}, 100)
				envy.ErrorChannels[functionName] = make(chan error, 100)
				envy.DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				envy.DefaultErrors[functionName] = fmt.Errorf("ruh roh")

				obj := envy.GetObject(functionName)
				g.Expect(obj).To(BeAssignableToTypeOf(struct{ Foo string }{}))
				g.Expect(obj.(struct{ Foo string })).To(Equal(struct{ Foo string }{Foo: "bar"}))

				err := envy.GetError(functionName)
				g.Expect(err).To(MatchError("ruh roh"))
			},
		},
		{
			testCase: "returns pushed object and error",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"

				envy.ObjectChannels[functionName] = make(chan interface{}, 100)
				envy.ErrorChannels[functionName] = make(chan error, 100)
				envy.DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				envy.DefaultErrors[functionName] = nil

				err := envy.AddObjectReturns(functionName, struct{ Foo string }{Foo: "foo"})
				g.Expect(err).NotTo(HaveOccurred())

				err = envy.AddErrorReturns(functionName, fmt.Errorf("some error"))
				g.Expect(err).NotTo(HaveOccurred())

				obj := envy.GetObject(functionName)
				g.Expect(obj).To(BeAssignableToTypeOf(struct{ Foo string }{}))
				g.Expect(obj.(struct{ Foo string })).To(Equal(struct{ Foo string }{Foo: "foo"}))

				err = envy.GetError(functionName)
				g.Expect(err).To(MatchError("some error"))
			},
		},
		{
			testCase: "error for function not found",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				functionName := "foobar"

				envy.ObjectChannels[functionName] = make(chan interface{}, 100)
				envy.ErrorChannels[functionName] = make(chan error, 100)
				envy.DefaultObjects[functionName] = struct{ Foo string }{Foo: "bar"}
				envy.DefaultErrors[functionName] = nil

				err := envy.AddObjectReturns("notafunction", struct{ Foo string }{Foo: "foo"})
				g.Expect(err).To(HaveOccurred())

				err = envy.AddErrorReturns("notafunction", fmt.Errorf("some error"))
				g.Expect(err).To(HaveOccurred())

			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
