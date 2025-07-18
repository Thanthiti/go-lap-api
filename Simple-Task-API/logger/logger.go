package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger{
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}