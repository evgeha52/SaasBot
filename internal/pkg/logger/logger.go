package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(level string) {
	cfg := zap.NewProductionConfig()
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	}
	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
