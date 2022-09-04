package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer   []io.Writer
	LogLevel []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, writer := range hook.Writer {
		writer.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevel
}

type Logger struct {
	*logrus.Entry
}

var e *logrus.Entry

func GetLogger() *Logger {
	return &Logger{e}
}

func GetLoggerWithField(k string, v any) *Logger {
	return &Logger{logrus.WithField(k, v)}
}

func init() {
	l := logrus.New()

	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
		DisableColors: false,
	}
	if err := os.Mkdir("logs", 0644); err != nil && !os.IsExist(err) {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 8640)
	if err != nil {
		panic(err)
	}
	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:   []io.Writer{allFile, os.Stdout},
		LogLevel: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
