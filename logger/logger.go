package logger

import (
	"github.com/sirupsen/logrus"

	"fmt"

	"github.com/qa-dev/jsonwire-grid/config"
)

// Init - initialisation of logger.
func Init(cfgLogger config.Logger) error {
	level, err := logrus.ParseLevel(cfgLogger.Level)
	if err != nil {
		return fmt.Errorf("parse log level, %v", err)
	}
	logrus.Infof("Set log level to: %v", level)
	logrus.SetLevel(level)
	return nil
}
