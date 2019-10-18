package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func initLogger() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	Logger = zap.New(core)
}

func logAndPrint(msg string) {
	fmt.Println(msg)
	Logger.Info(msg)
}
