package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"
)

// ContextKey is the type used in setting
type ContextKey string

const (
	TRACEID = "traceID"
	CALLER  = "caller"
	TS      = "ts"

	KeyTraceID ContextKey = TRACEID
)

// Logger provides an interface for implementing
// a MDC logger
type Logger interface {
	Error(context.Context, error)
	Errorf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Panic(context.Context, error)
	Panicf(context.Context, string, ...interface{})
}

// Logger is an logger implementation based on
// uber zapp
type zapLogger struct {
	Logger      *zap.Logger
	Application string
}

// NewZapLogger returns a new logger backed by Zap
func NewZapLogger(application, env string) Logger {

	rawJSONTpl := `{
					"level": "debug",
					"encoding": "json",
					"outputPaths": ["stdout"],
					"errorOutputPaths": ["stderr"],
					"initialFields": {"application": "%s", "environment": "%s"},
					"encoderConfig" : {
						"messageKey": "message",
						"levelKey": "level",
						"levelEncoder": "lowercase"
					}
				}`

	rawJSON := fmt.Sprintf(rawJSONTpl, application, env)

	var cfg zap.Config

	if err := json.Unmarshal([]byte(rawJSON), &cfg); err != nil {
		panic(err)
	}

	l, e := cfg.Build()

	if e != nil {
		panic(e)
	}

	return &zapLogger{
		Logger:      l,
		Application: application,
	}
}

// Infof logs a message with criticality INFO
func (logger *zapLogger) Infof(ctx context.Context, format string, a ...interface{}) {
	defer logger.Logger.Sync()

	logger.Logger.Info(fmt.Sprintf(format, a...), buildFields(ctx)...)
}

// Error logs an error with criticality ERROR
func (logger *zapLogger) Error(ctx context.Context, e error) {
	defer logger.Logger.Sync()

	logger.Logger.Error(e.Error(), buildFields(ctx)...)
}

// Errorf logs an error message with criticality ERROR
func (logger *zapLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	defer logger.Logger.Sync()

	logger.Logger.Error(fmt.Sprintf(format, a...), buildFields(ctx)...)
}

// Panic logs an error with criticality PANIC
func (logger *zapLogger) Panic(ctx context.Context, e error) {
	defer logger.Logger.Sync()

	logger.Logger.Panic(e.Error(), buildFields(ctx)...)
}

// Panicf logs an error message with criticality PANIC
func (logger *zapLogger) Panicf(ctx context.Context, format string, a ...interface{}) {
	defer logger.Logger.Sync()

	logger.Logger.Panic(fmt.Sprintf(format, a...), buildFields(ctx)...)
}

// Warnf logs a message with criticality WARN
func (logger *zapLogger) Warnf(ctx context.Context, format string, a ...interface{}) {
	defer logger.Logger.Sync()

	logger.Logger.Warn(fmt.Sprintf(format, a...), buildFields(ctx)...)
}

func buildFields(ctx context.Context) []zap.Field {
	out := make([]zap.Field, 0, 3)

	// Set trace id
	if traceID := ctx.Value(KeyTraceID); traceID != nil {
		out = append(out, zap.String(TRACEID, traceID.(string)))
	}

	// extract and set caller
	_, file, line, ok := runtime.Caller(2)

	if ok {
		out = append(out, zap.String(CALLER, fmt.Sprintf("%s:%d", file, line)))
	}

	// set timestamp

	ts := time.Now()

	out = append(out, zap.String(TS, fmt.Sprintf("%d", ts.UnixNano())))

	return out
}
