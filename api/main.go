package main

import (
	"github.com/bitcrshr/envmgr/api/grpc"
	"github.com/bitcrshr/envmgr/api/shared"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	defer shared.Logger().Sync()

	if err := godotenv.Load(); err != nil {
		shared.Logger().Fatal("failed to load env", zap.Error((err)))
	}

	grpc.Serve()
}
