package platform

import (
	"os"
	"os/user"
	"path"

	"github.com/lorentzforces/selfman/internal/run"
)

func ResolveXdgConfigDir() string {
	xdgEnvPath := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgEnvPath) > 0 {
		return xdgEnvPath
	}

	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".config")
}

func ResolveXdgDataDir() string {
	xdgEnvPath := os.Getenv("XDG_DATA_HOME")
	if len(xdgEnvPath) > 0 {
		return xdgEnvPath
	}

	usr, err := user.Current()
	run.AssertNoErr(err)
	return path.Join(usr.HomeDir, ".local", "share")
}

func ResolveRepoPathForApp(name string) string {
	basePath := ResolveXdgDataDir()
	return path.Join(basePath, "selfman", "repos", name)
}
