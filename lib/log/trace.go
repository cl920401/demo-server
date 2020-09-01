package log

import (
	"demo-server/lib/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TraceLog(path string) *zap.Logger {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:  path,
		MaxSize:   config.Get("cm.log.size").Int(1024), //MB
		LocalTime: true,
		Compress:  config.Get("cm.log.compress").Bool(true),
	})

	var traceLogLevel = zap.NewAtomicLevel()
	traceLogLevel.SetLevel(zap.InfoLevel)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := NewRawEncoder(encoderConfig)
	fileWriter := zapcore.NewCore(
		fileEncoder,
		file,
		traceLogLevel,
	)

	core := zapcore.NewTee(fileWriter)
	return zap.New(core)
}
