package loggerinterface

//go:generate mockgen -destination=mocks/logger_mock.go -package=mocks authstore/internal/common/loggerinterface Logger
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
