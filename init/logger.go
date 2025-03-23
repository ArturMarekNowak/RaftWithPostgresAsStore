package init

import "go.uber.org/zap"

func ConfigureLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return sugar
}
