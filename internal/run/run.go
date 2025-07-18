package run

import (
	"fmt"
	"io"
	"os"

	"github.com/lorentzforces/fresh-err/fresherr"
	"gopkg.in/yaml.v3"
)

func FailOut(msg string) {
	fmt.Fprintln(os.Stderr, ErrMsg(msg))
	os.Exit(1)
}

func FailOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, ErrMsg(err.Error()))
		os.Exit(1)
	}
}

func ErrMsg(msg string) string {
	return fresherr.GetFresh() + ": " + msg
}

func Assert(condition bool, more any) {
	if condition { return }
	panic(fmt.Sprintf("Assertion violated: %s", more))
}

func AssertNoErr(err error) {
	if err == nil { return }
	panic(fmt.Sprintf("Assertion violated, error encountered: %s ", err.Error()))
}

func AssertNoErrReason(err error, reason string) {
	if err == nil { return }
	panic(fmt.Sprintf("Assertion violated, error encountered (%s): %s ", reason, err.Error()))
}

func Coalesce[T any](a, b *T) *T {
	if a == nil { return b }
	return a
}

func CoalesceString(a, b string) string {
	if len(a) == 0 { return b }
	return a
}

// Returns a pointer to a passed string, for when you want to put a string into a *string field in
// a literal.
//
// i.e. myObj{ str: run.StrPtr("a string literal") }
func StrPtr(str string) *string {
	return &str
}

var ErrNotImplemented = fmt.Errorf("Not yet implemented")

func VerifyDirExists(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GetStrictDecoder(source io.Reader) *yaml.Decoder {
	decoder := yaml.NewDecoder(source)
	decoder.KnownFields(true)
	return decoder
}
