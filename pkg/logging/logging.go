package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerMap sync.Map

func initLogger(name string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:  filepath.Join("logs", fmt.Sprintf("%s.log", name)),
		MaxSize:   1024,
		MaxAge:    365,
		Compress:  true,
		LocalTime: true,
	}

	fileWriter := zapcore.AddSync(&hook)
	errorPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l == zap.ErrorLevel })

	consoleDebugging := zapcore.Lock(os.Stdout)

	fileConfig := zap.NewProductionEncoderConfig()
	fileConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(fileConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)

	stack := zap.AddStacktrace(errorPriority)
	caller := zap.AddCaller()

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, zap.NewAtomicLevelAt(zapcore.DebugLevel)),
		zapcore.NewCore(consoleEncoder, consoleDebugging, zap.NewAtomicLevelAt(zapcore.DebugLevel)),
	)

	logger := zap.New(core, stack, caller)

	return logger
}

func GetLogger(name string) *zap.Logger {
	v, ok := loggerMap.Load(name)
	if !ok {
		logger := initLogger(name)
		loggerMap.Store(name, logger)

		return logger
	}

	return v.(*zap.Logger)
}

func GetSugaredLogger(name string) *zap.SugaredLogger {
	v, ok := loggerMap.Load(name)
	if !ok {
		logger := initLogger(name).Sugar()
		loggerMap.Store(name, logger)

		return logger
	}

	return v.(*zap.SugaredLogger)
}
