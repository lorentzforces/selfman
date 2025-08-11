package data

import (
	"os"
	"os/user"
	"path"

	"github.com/lorentzforces/selfman/internal/run"
)

func resolveXdgConfigDir() string {
	xdgEnvPath := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgEnvPath) > 0 {
		return xdgEnvPath
	}

	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".config")
}

func resolveXdgDataDir() string {
	xdgEnvPath := os.Getenv("XDG_DATA_HOME")
	if len(xdgEnvPath) > 0 {
		return xdgEnvPath
	}

	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".local", "share")
}

func resolveXdgBinDir() string {
	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".local", "bin")
}

func resolveUserLibDir() string {
	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".local", "lib")
}
