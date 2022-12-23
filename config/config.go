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
		Level string
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
	overwriteStrIfEnv(&c.Storage.Path, "GMON_DATA_DIR")
	c.Storage.Path = os.ExpandEnv(c.Storage.Path)
}
