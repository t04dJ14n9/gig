package main

import (
	"bytes"
	logrus "github.com/sirupsen/logrus"
)

// LogrusLevelParse tests parsing a log level string
func LogrusLevelParse(level string) string {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return "error"
	}
	return l.String()
}

// LogrusNewLogger tests creating a new logger and setting format
func LogrusNewLogger() string {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Info("hello")
	return buf.String()
}

// LogrusGetLevel tests getting the current log level
func LogrusGetLevel() string {
	return logrus.GetLevel().String()
}
