package envy

import (
	"errors"
	"os"
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
)

type testRunner struct {
	testCase string
	runner   func(tt *testing.T)
}

func TestEnv(t *testing.T) {
	tests := []testRunner{
		{
			testCase: "sets and unsets variables",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				err := os.Setenv("TEST_VAR1", "foo")
				g.Expect(err).NotTo(HaveOccurred())

				err = os.Setenv("TEST_VAR2", "biz")
				g.Expect(err).NotTo(HaveOccurred())

				err = os.Unsetenv("TEST_VAR3")
				g.Expect(err).NotTo(HaveOccurred())

				env := MockEnv{}
				err = env.Load(
					map[string]string{
						"TEST_VAR1": "bar",
						"TEST_VAR2": "", // should unset empty variables
						"TEST_VAR3": "baz",
					},
				)
				g.Expect(err).NotTo(HaveOccurred())

				var1, exists := os.LookupEnv("TEST_VAR1")
				g.Expect(exists).To(BeTrue())
				g.Expect(var1).To(Equal("bar"))

				_, exists = os.LookupEnv("TEST_VAR2")
				g.Expect(exists).To(BeFalse())

				var3, exists := os.LookupEnv("TEST_VAR3")
				g.Expect(exists).To(BeTrue())
				g.Expect(var3).To(Equal("baz"))

				env.Restore()

				var1, exists = os.LookupEnv("TEST_VAR1")
				g.Expect(exists).To(BeTrue())
				g.Expect(var1).To(Equal("foo"))

				var2, exists := os.LookupEnv("TEST_VAR2")
				g.Expect(exists).To(BeTrue())
				g.Expect(var2).To(Equal("biz"))

				_, exists = os.LookupEnv("TEST_VAR3")
				g.Expect(exists).To(BeFalse())

			},
		},
		{
			testCase: "panics for restoring unloaded environment",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := MockEnv{}

				g.Expect(env.Restore).To(PanicWith(errors.New("environment not loaded")))

			},
		},
		{
			testCase: "panics for loading a loaded environment",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := MockEnv{}
				env.Load(map[string]string{})

				g.Expect(env.Load(map[string]string{})).To(MatchError("mock environment is already loaded"))

			},
		},
		{
			testCase: "panics for loading empty variable key",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := MockEnv{}

				g.Expect(env.Load(map[string]string{"": "bad key"})).To(MatchError(syscall.EINVAL))

			},
		},
		{
			testCase: "panics for restoring empty variable key",
			runner: func(tt *testing.T) {
				g := NewGomegaWithT(tt)

				env := MockEnv{isLoaded: true, existing: map[string]string{"": "bad key"}}

				g.Expect(env.Restore).To(PanicWith(os.NewSyscallError("setenv", syscall.EINVAL)))

			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.testCase, test.runner)
	}
}
