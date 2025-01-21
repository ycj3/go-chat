package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	// Set log level to Debug
	logrus.SetLevel(logrus.DebugLevel)
}
