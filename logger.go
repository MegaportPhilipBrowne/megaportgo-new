package megaport

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

const (
	TraceLevel = slog.Level(-8)
	DebugLevel = slog.Level(-4)
	InfoLevel  = slog.Level(0)
	WarnLevel  = slog.Level(4)
	ErrorLevel = slog.Level(8)
	FatalLevel = slog.Level(12)
	Off        = slog.Level(16)
)

var LevelNames = map[slog.Leveler]string{
	TraceLevel: "TRACE",
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
	Off:        "NONE",
}

func StringToLogLevel(level string) slog.Level {
	switch level {
	case "TRACE":
		return TraceLevel
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return Off
	}
}

type DefaultLogger struct {
	level slog.Level
}

func NewDefaultLogger() *DefaultLogger {
	d := DefaultLogger{level: DebugLevel}
	return &d
}

func (d *DefaultLogger) SetLevel(l slog.Level) {
	d.level = l
}

func (d *DefaultLogger) log(level slog.Level, msg string, args ...interface{}) {
	if level >= d.level {
		opts := &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					levelLabel, exists := LevelNames[level]
					if !exists {
						levelLabel = level.String()
					}

					a.Value = slog.StringValue(levelLabel)
				}
				return a
			},
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		ctx := context.Background()
		logger.Log(ctx, level, msg, args...)
	}
}

// Emit the message and args at DEBUG level
func (d *DefaultLogger) Debug(msg string, args ...interface{}) {
	d.log(DebugLevel, msg, args...)
}

// Emit the message and args at TRACE level
func (d *DefaultLogger) Trace(msg string, args ...interface{}) {
	d.log(TraceLevel, msg, args...)
}

// Emit the message and args at INFO level
func (d *DefaultLogger) Info(msg string, args ...interface{}) {
	d.log(InfoLevel, msg, args...)
}

// Emit the message and args at WARN level
func (d *DefaultLogger) Warn(msg string, args ...interface{}) {
	d.log(WarnLevel, msg, args...)
}

// Emit the message and args at ERROR level
func (d *DefaultLogger) Error(msg string, args ...interface{}) {
	d.log(ErrorLevel, msg, args...)
}

// Emit the message and args at FATAL level
func (d *DefaultLogger) Fatal(msg string, args ...interface{}) {
	d.log(FatalLevel, msg, args...)
}
