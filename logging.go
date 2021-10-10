package main

import (
	dapnet "github.com/michaelpeterswa/godapnet"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
	logger.Info("DAPNetSendARRLNews is initializing...")
	dapnet.InitializeLogger(logger)
}
