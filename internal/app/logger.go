package app

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func InitLogger(level, filepath, filename string) *logrus.Logger {
	l := map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
	}

	log := logrus.New()
	log.SetOutput(os.Stdout)

	if level != "" {
		log.SetLevel(l[level])

		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			log.WithField("path", filepath).Info(fmt.Errorf("Failed to create floder: %w", err))
		}
		p := path.Join(filepath, filename)

		file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using stdout.")
		}

		log.Info("Logger has initialized.")
	}

	return log
}
