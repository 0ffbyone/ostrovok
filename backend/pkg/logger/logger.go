package logger

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(args ...any)
	Info(args ...any)
	Warn(args ...any)
}

type logger struct{ *zap.SugaredLogger }

func NewLogger() (Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		log.Panic("failed to create logger", err)
	}
	defer zapLogger.Sync()

	sugar := zapLogger.Sugar()
	if sugar == nil {
		return nil, errors.New("sugared logger failed to start")
	}

	return logger{sugar}, nil
}

func (l logger) Error(args ...any) {
	args = createArgs(args...)
	l.SugaredLogger.Error(args)
}

func (l logger) Info(args ...any) {
	args = createArgs(args...)
	l.SugaredLogger.Info(args)
}

func (l logger) Warn(args ...any) {
	args = createArgs(args...)
	l.SugaredLogger.Warn(args)
}

func isStruct(i any) bool {
	t := reflect.TypeOf(i)

	return t.Kind() == reflect.Struct
}

func createArgs(args ...any) []any {
	for i := range args {
		curr := args[i]
		if isStruct(args[i]) {
			typeOf := reflect.TypeOf(curr)
			name := typeOf.Name()
			args[i] = fmt.Sprintf("%s:%+v", name, curr)
		}
	}

	return args
}
