package main

import (
	"context"
	"net/http"

	"https://github.com/IminaDanzi/Docker-Project/dat"
	"https://github.com/IminaDanzi/Docker-Project/lib"
	"https://github.com/IminaDanzi/Docker-Project/logger"
	"go.uber.org/zap"
)

func main() {
	logger := logger.InitLogger()
	ctx := context.Background()

	db, err := dat.InitDB(ctx, logger)
	if err != nil {
		logger.Fatal("failed to init DB", zap.Error(err))
	}
	defer dat.ExitDb(logger)

	_, err = db.Exec(ctx, dat.GetSchemaSQL())
	if err != nil {
		logger.Fatal("failed to execute SQL schema", zap.Error(err))
	}

	address := ":8080"
	server := http.Server{
		Addr:    address,
		Handler: lib.NewRouter(),
	}

	logger.Info("Listening...", zap.String("address", address))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("ListenAndServe failed", zap.Error(err))
	}
}
