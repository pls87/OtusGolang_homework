package logger

import (
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"os"

	"github.com/sirupsen/logrus"
)

func New(cfg config.LoggerConf) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	if lvl, err := logrus.ParseLevel(cfg.Level); err != nil {
		log.Level = lvl
	}
	return log
}
