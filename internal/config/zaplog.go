package config

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var syncOnceLog sync.Once
var logger *zap.Logger

type EncoderOutput int32

const (
	JSON EncoderOutput = iota
	Default
)

func (o EncoderOutput) String() string {
	switch o {
	case JSON:
		return "JSON"
	case Default:
		return "DEFAULT"
	}
	return "unknown"
}

func ProvideLogger(config *Config) (*zap.Logger, error) {
	var err error 
	syncOnceLog.Do(func() {
		err = os.MkdirAll(config.LogDir, 0755)
		if err != nil {
			return
		}
		encoderJSON := encoder(JSON)
		encoderDefault := encoder(Default)
		core := zapcore.NewTee(
			zapcore.NewCore(encoderDefault, zapcore.AddSync(os.Stdout), zap.DebugLevel),
			zapcore.NewCore(encoderJSON, logWriter(config.LogDir, "info.log"), zap.InfoLevel),
			zapcore.NewCore(encoderJSON, logWriter(config.LogDir, "error.log"), zap.ErrorLevel),
		)
		logger = zap.New(core)
		defer logger.Sync()
	})
	return logger, err
}

func encoder(o EncoderOutput) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	// https://github.com/uber-go/zap/blob/master/zapcore/encoder.go
	// config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05")

	switch o {
	case JSON:
		return zapcore.NewJSONEncoder(config)
	case Default:
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

func logWriter(logDir string, filename string) zapcore.WriteSyncer {
	fullPath := path.Join(logDir, filename)
	lumberjack := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    1000, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	}
	return zapcore.AddSync(lumberjack)
}

func TimeTrack(start time.Time, fn string, name string) {
	elapsed := time.Since(start)
	logger.Info(
		fmt.Sprintf("[ECHO|%s] %v | %s | %s", fn, time.Now().Format("2006/01/02 - 15:04:05"),
			name,
			elapsed,
		),
	)
}