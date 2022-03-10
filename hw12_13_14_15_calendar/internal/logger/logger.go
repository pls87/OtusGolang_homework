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

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	return log
}
