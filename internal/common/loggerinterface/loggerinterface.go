package loggerinterface

type Logger interface {
	Info(...any)
	Infof(string, ...any)
	Error(...any)
	Errorf(string, ...any)
	Warn(...any)
	Warnf(string, ...any)
	Fatal(...any)
	Fatalf(string, ...any)
}
