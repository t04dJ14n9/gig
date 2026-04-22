package main

import "go.uber.org/zap"

func ZapNewLogger() string {
	logger, _ := zap.NewProduction()
	if logger == nil {
		return "nil"
	}
	return "ok"
}

func ZapSugarLogger() string {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	if sugar == nil {
		return "nil"
	}
	return "ok"
}

func ZapNamedLogger() string {
	logger, _ := zap.NewProduction()
	named := logger.Named("test")
	if named == nil {
		return "nil"
	}
	return "ok"
}
