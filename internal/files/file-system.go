package files

import "github.com/lorentzforces/selfman/internal/config"

// TODO/WIP: I need to figure out the proper interface for stuff that is going to live on-disk.
// Just using the file system is all well and good, but I want to actually write _tests_ for the
// damn thing.

type FileSystem struct {
}

func (self FileSystem) AppSourceExists(app config.AppConfig) bool {
	return false
}
