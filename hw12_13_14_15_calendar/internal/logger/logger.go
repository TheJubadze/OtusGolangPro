package logger

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Setup(logLevel string) {
	log = logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Error parsing log level: %s", err)
	}
	log.SetLevel(level)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func Log() *logrus.Logger {
	return log
}
