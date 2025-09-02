package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(isDev bool) {

	if isDev {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}
}
