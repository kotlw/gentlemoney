package config

import "time"

func Default() *Config {
	const defaultPath = "$HOME/.gentlemoney"

	c := &Config{
		App: App{
			Name:    "gentlemoney",
			Version: "v0.1",
		},
		Logger: Logger{
			Path:     defaultPath + "/logs",
			Filename: "log_" + time.Now().Format("2006-01-02T15:04:05 -07:00:00"),
			Level:    "",
		},
		Storage: Storage{
			Path:     defaultPath,
			Filename: "data.sqlite3",
		},
	}

	postprocess(c)

	return c
}
