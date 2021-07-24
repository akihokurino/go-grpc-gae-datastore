package adapter

import "context"

type CompositeLogger interface {
	For(l LogLevel) ContextualLogger
	Debug() ContextualLogger
	Info() ContextualLogger
	Warn() ContextualLogger
	Error() ContextualLogger
}

type Logger interface {
	Print(v interface{})
	Printf(format string, args ...interface{})
}

type ContextualLogger interface {
	Logger
	With(ctx context.Context) Logger
}

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

const (
	LogLevelDebug LogLevel = iota // 開発時にのみ表示したいもの
	LogLevelInfo                  // 問題ではないが後からトラッキングしたいもの
	LogLevelWarn                  // エラーではないが頻発する場合関心をもつべきもの
	LogLevelError                 // 発生したら原則として開発者が即時対応をするべきで、ゼロに近づけることを目指すべきもの
)
