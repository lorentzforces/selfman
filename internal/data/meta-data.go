package data

type Meta struct {
	CurrentVersion string `yaml:"current-version"`
	// TODO: version history
}

func MetaFileNameForApp(appName string) string {
	return appName + ".meta.yaml"
}
