package main

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func InitLog() {
	if Log != nil {
		return
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}

func SyncLog() {
	Log.Sync()
}
