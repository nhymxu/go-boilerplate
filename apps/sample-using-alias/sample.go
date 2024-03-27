package sample_using_alias

import "go.uber.org/zap"

func Run() {
	logger := zap.S()

	logger.Info("Sample app using alias project name")
}
