package config

import (
	"os"
)

type (
	// Config - main config structure.
	Config struct {
		App     App
		Logger  Logger
		Storage Storage
	}

	// App - application config.
	App struct {
		Name    string
		Version string
	}

	// Logger - logger config.
	Logger struct {
		Path     string
		Filename string
		Level    string
	}

	// Storage - storage config.
	Storage struct {
		Path     string
		Filename string
	}
)

func overwriteStrIfEnv(targetValue *string, envKey string) {
	if envValue := os.Getenv(envKey); envValue != "" {
		*targetValue = envValue
	}
}

func postprocess(c *Config) {
	// Logger path
	overwriteStrIfEnv(&c.Logger.Path, "GMON_DATA_DIR")
	c.Logger.Path = os.ExpandEnv(c.Logger.Path) + "/logs"
	overwriteStrIfEnv(&c.Logger.Path, "GMON_LOG_DIR")
	c.Logger.Path = os.ExpandEnv(c.Logger.Path)

	// Logger level
	overwriteStrIfEnv(&c.Logger.Level, "GMON_LOG")
	c.Logger.Level = os.ExpandEnv(c.Logger.Level)

	// Storage path
	overwriteStrIfEnv(&c.Storage.Path, "GMON_DATA_DIR")
	c.Storage.Path = os.ExpandEnv(c.Storage.Path)
}
