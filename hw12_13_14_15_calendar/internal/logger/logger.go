package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	if lvl, err := logrus.ParseLevel(level); err != nil {
		log.Level = lvl
	}
	return log
}
