package gd

func LogClose() {
}

func Debug(arg0 interface{}, args ...interface{}) {
}

func Crash(args ...interface{}) {
}

func Crashf(format string, args ...interface{}) {
}

func Exit(args ...interface{}) {
}

func Exitf(format string, args ...interface{}) {
}

func Stderr(args ...interface{}) {
}

func Stderrf(format string, args ...interface{}) {
}

func Stdout(args ...interface{}) {
}

func Stdoutf(format string, args ...interface{}) {
}

func GetLevel() string {
	return ""
}

func SetLevel(lvl int) {
}

func Log(lvl string, source, message string) {
}

func Logf(lvl string, format string, args ...interface{}) {
}

func Logc(lvl string, closure func() string) {
}

func Finest(arg0 interface{}, args ...interface{}) {
}

func Fine(arg0 interface{}, args ...interface{}) {
}

func DebugT(tag string, arg0 interface{}, args ...interface{}) {
}

func Trace(arg0 interface{}, args ...interface{}) {
}

func TraceT(tag string, arg0 interface{}, args ...interface{}) {
}

func Info(arg0 interface{}, args ...interface{}) {
}

func InfoT(tag string, arg0 interface{}, args ...interface{}) {
}

func Warn(arg0 interface{}, args ...interface{}) {
}

func WarnT(tag string, arg0 interface{}, args ...interface{}) {
}

func Error(arg0 interface{}, args ...interface{}) {
}

func ErrorT(tag string, arg0 interface{}, args ...interface{}) {
}

func Critical(arg0 interface{}, args ...interface{}) {
}

func CriticalT(tag string, arg0 interface{}, args ...interface{}) {
}
