package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"gae-go-recruiting-server/adapter"
)

func NewLoggerWithMinLevel(minLevel adapter.LogLevel) adapter.CompositeLogger {
	loggerFor := func(level adapter.LogLevel, logger *log.Logger) adapter.ContextualLogger {
		return ContextualLogEmitter(func(ctx context.Context, v interface{}) {
			logger.Printf("[%s] %+v", level.String(), v)
		})
	}

	return &StaticCompositeLogger{
		Loggers: [5]adapter.ContextualLogger{
			loggerFor(adapter.LogLevelDebug, StdoutSingletonLogger),
			loggerFor(adapter.LogLevelInfo, StdoutSingletonLogger),
			loggerFor(adapter.LogLevelWarn, StderrSingletonLogger),
			loggerFor(adapter.LogLevelError, StderrSingletonLogger),
		},
		MinLevel: minLevel,
	}
}

var NopContextualLogger = ContextualLogEmitter(func(ctx context.Context, v interface{}) {})

type logEmitter func(v interface{})

func (e logEmitter) Print(v interface{}) {
	e(v)
}

func (e logEmitter) Printf(format string, args ...interface{}) {
	e(fmt.Sprintf(format, args...))
}

type StaticCompositeLogger struct {
	Loggers  [5]adapter.ContextualLogger
	MinLevel adapter.LogLevel
}

type ContextualLogEmitter func(ctx context.Context, v interface{})

func (e ContextualLogEmitter) Print(v interface{}) {
	e(nil, v)
}

func (e ContextualLogEmitter) Printf(format string, args ...interface{}) {
	e(nil, fmt.Sprintf(format, args...))
}

func (e ContextualLogEmitter) With(ctx context.Context) adapter.Logger {
	return logEmitter(func(v interface{}) { e(ctx, v) })
}

func (cl *StaticCompositeLogger) For(l adapter.LogLevel) adapter.ContextualLogger {
	if l < cl.MinLevel {
		return NopContextualLogger
	}
	return cl.Loggers[l]
}

func (cl *StaticCompositeLogger) Debug() adapter.ContextualLogger {
	return cl.For(adapter.LogLevelDebug)
}

func (cl *StaticCompositeLogger) Info() adapter.ContextualLogger {
	return cl.For(adapter.LogLevelInfo)
}

func (cl *StaticCompositeLogger) Warn() adapter.ContextualLogger {
	return cl.For(adapter.LogLevelWarn)
}

func (cl *StaticCompositeLogger) Error() adapter.ContextualLogger {
	return cl.For(adapter.LogLevelError)
}

var (
	StdoutSingletonLogger = log.New(os.Stdout, "", 0)
	StderrSingletonLogger = log.New(os.Stderr, "", 0)
)
