package logger

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Setup(logLevel string) {
	Log = logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Log.Fatalf("Error parsing Log level: %s", err)
	}
	Log.SetLevel(level)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		DisableColors:   false,
		DisableQuote:    true,
	})
}
