package logger

import (
	"github.com/Sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/config"
)

func Init(logger config.Logger) {
	logrus.SetLevel(logger.Level)
}
