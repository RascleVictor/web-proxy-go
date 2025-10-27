package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(level string) {
	var err error
	if level == "debug" {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
}
