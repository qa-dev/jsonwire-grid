package logger

import (
	"github.com/Sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/config"
	"fmt"
)

func Init(logger config.Logger) error {
	level, err := logrus.ParseLevel(logger.Level)
	if err != nil {
		return fmt.Errorf("Parse log level, %v", err)
	}
	logrus.Infof("Set log level to: %v", level)
	logrus.SetLevel(level)
	return nil
}
