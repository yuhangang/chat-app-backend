package log

import (
	"github.com/sirupsen/logrus"
)

var (
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
)

func SetupLoggers() {

	InfoLogger = logrus.New()
	ErrorLogger = logrus.New()

	// Set Logrus log level and formatter
	InfoLogger.SetLevel(logrus.InfoLevel)

	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
