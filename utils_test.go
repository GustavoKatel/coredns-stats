package stats_test

import (
	"log"

	stats "github.com/GustavoKatel/coredns-stats"
)

type logger struct{}

// Debug implements stats.Logger.
func (l *logger) Debug(v ...interface{}) {
	log.Print(v...)
}

// Debugf implements stats.Logger.
func (l *logger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Error implements stats.Logger.
func (l *logger) Error(v ...interface{}) {
	log.Print(v...)
}

// Errorf implements stats.Logger.
func (l *logger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Fatal implements stats.Logger.
func (l *logger) Fatal(v ...interface{}) {
	log.Print(v...)
}

// Fatalf implements stats.Logger.
func (l *logger) Fatalf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Info implements stats.Logger.
func (l *logger) Info(v ...interface{}) {
	log.Print(v...)
}

// Infof implements stats.Logger.
func (l *logger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Warning implements stats.Logger.
func (l *logger) Warning(v ...interface{}) {
	log.Print(v...)
}

// Warningf implements stats.Logger.
func (l *logger) Warningf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

var _ stats.Logger = (*logger)(nil)
