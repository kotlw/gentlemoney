package config

func Default() *Config {
	c := &Config{
		App: App{
			Name:    "gentlemoney",
			Version: "v0.1",
		},
		Logger: Logger{
			Level: "debug",
		},
		Storage: Storage{
			Path:     "$HOME/.gentlemoney",
			Filename: "data.sqlite3",
		},
	}

	postprocess(c)

	return c
}
