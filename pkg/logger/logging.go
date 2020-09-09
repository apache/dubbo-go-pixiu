package logger

// Info
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Debug
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Infof
func Infof(fmt string, args ...interface{}) {
	logger.Infof(fmt, args...)
}

// Warnf
func Warnf(fmt string, args ...interface{}) {
	logger.Warnf(fmt, args...)
}

// Errorf
func Errorf(fmt string, args ...interface{}) {
	logger.Errorf(fmt, args...)
}

// Debugf
func Debugf(fmt string, args ...interface{}) {
	logger.Debugf(fmt, args...)
}
