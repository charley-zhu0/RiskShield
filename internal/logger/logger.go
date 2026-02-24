package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func init() {
	logger = zap.NewNop()
}

func Init(logDir string, maxBackups int) error {
	absPath, err := filepath.Abs(logDir)
	if err != nil {
		return fmt.Errorf("invalid log directory: %w", err)
	}
	if filepath.Clean(absPath) != absPath {
		return fmt.Errorf("log directory contains invalid path")
	}

	if err := os.MkdirAll(absPath, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	logDir = absPath

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	infoWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "info.log"),
		MaxSize:    100,
		MaxBackups: maxBackups,
	}

	errorWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "error.log"),
		MaxSize:    100,
		MaxBackups: maxBackups,
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), infoLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Sync() error {
	return logger.Sync()
}
