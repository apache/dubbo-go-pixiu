package logger

func Info(args ...interface{}) {
	logger.Info(args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Infof(fmt string, args ...interface{}) {
	logger.Infof(fmt, args...)
}
func Warnf(fmt string, args ...interface{}) {
	logger.Warnf(fmt, args...)
}
func Errorf(fmt string, args ...interface{}) {
	logger.Errorf(fmt, args...)
}
func Debugf(fmt string, args ...interface{}) {
	logger.Debugf(fmt, args...)
}
