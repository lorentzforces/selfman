package run

import (
	"fmt"
	"os"
	"os/user"
	"path"
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
	return "ERROR: " + msg
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
	if a == nil {
		return b
	}
	return a
}

func ResolveXdgConfigDir() string {
	xdgEnvPath := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgEnvPath) > 0 {
		return xdgEnvPath
	}

	usr, err := user.Current()
	AssertNoErr(err) // unsure how it's possible for this to fail
	return path.Join(usr.HomeDir, ".config")
}
