package logger

import "github.com/golang/glog"

const debugVerosity = 10

// Logger is interface for logging service
type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// Log basic log
type Log struct {
}

// Infof logs at info level
func (l Log) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

// Errorf logs at error level
func (l Log) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

// Debugf logs at debug level
func (l Log) Debugf(format string, args ...interface{}) {
	glog.V(debugVerosity).Infof(format, args...)
}
