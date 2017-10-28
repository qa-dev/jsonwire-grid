package logger

import (
	"github.com/Sirupsen/logrus"

	"fmt"
	"github.com/qa-dev/jsonwire-grid/config"
)

// Init - initialisation of logger.
func Init(logger config.Logger) error {
	level, err := logrus.ParseLevel(logger.Level)
	if err != nil {
		return fmt.Errorf("Parse log level, %v", err)
	}
	logrus.Infof("Set log level to: %v", level)
	logrus.SetLevel(level)
	return nil
}
