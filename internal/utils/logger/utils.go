package logger

func Info(msg string, keysAndValues ...interface{}) {
	zapLogSugared.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	zapLogSugared.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	zapLogSugared.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	zapLogSugared.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	zapLogSugared.Fatalw(msg, keysAndValues...)
}

func Infoln(args ...interface{}) {
	zapLogSugared.Infoln(args...)
}

func Debugln(args ...interface{}) {
	zapLogSugared.Debugln(args...)
}

func Warnln(args ...interface{}) {
	zapLogSugared.Warnln(args...)
}

func Errorln(args ...interface{}) {
	zapLogSugared.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	zapLogSugared.Fatalln(args...)
}

func Infof(msg string, args ...interface{}) {
	zapLogSugared.Infof(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	zapLogSugared.Debugf(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	zapLogSugared.Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	zapLogSugared.Errorf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	zapLogSugared.Fatalf(msg, args...)
}
