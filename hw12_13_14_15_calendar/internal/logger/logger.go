package logger

import (
	"os"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/sirupsen/logrus"
)

func New(cfg configs.LoggerConf) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	if lvl, err := logrus.ParseLevel(cfg.Level); err != nil {
		log.Level = lvl
	}
	return log
}
