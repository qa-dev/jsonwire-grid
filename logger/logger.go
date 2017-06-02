package logger

import (
	"github.com/Sirupsen/logrus"

	"jsonwire-grid/config"
)

func Init(logger config.Logger) {
	logrus.SetLevel(logger.Level)
}
