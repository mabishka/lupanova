package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mabishka/lupanova/internal/compress"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/repository/connloader"
	"github.com/mabishka/lupanova/internal/repository/fileloader"
)

func main() {
	ctx, fnCancel := context.WithCancelCause(context.Background())
	defer fnCancel(errors.New("exit"))
	run(ctx)
}

func run(ctx context.Context) {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	logger.Log().Info("config",
		zap.String("serverAddress", config.GetServerAddress()),
		zap.String("baseAddress", config.GetBaseAddress()),
		zap.String("logLevel", config.GetLogLevel()),
		zap.String("fileName", config.GetFileName()),
		zap.String("connAddress", config.GetConnAddress()),
	)

	server := handler.New(config.GetBaseAddress())

	var loader model.StorageLoader
	var conn model.ConnLoader
	if config.GetConnAddress() != "" {
		loader = connloader.New(config.GetConnAddress())
		conn, _ = loader.(model.ConnLoader)
		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Info("conn not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("conn storage usage")
		}

	}

	if loader == nil && config.GetFileName() != "" {
		loader = fileloader.New(config.GetFileName())

		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Error("file not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("file storage usage")
		}
	}

	if loader == nil {
		logger.Log().Info("memory storage usage")
	}

	connServer := handler.NewConn(conn)

	router := chi.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(compress.WithCompress)

	router.Post("/", server.HandlerPostFull)
	router.Post("/api/shorten", server.HandlerPostFullJSON)
	router.Post("/api/shorten/batch", server.HandlerPostBatch)
	router.Get("/{id}", server.HandlerGetFull)
	router.Get("/ping", connServer.HandlerGetPing)

	go func() {
		if err := http.ListenAndServe(config.GetServerAddress(), router); err != nil {
			panic(err)
		}
	}()

	logger.Log().Info("listen port", zap.String("serverAddress", config.GetServerAddress()))

	<-ctx.Done()
	logger.Log().Info("exit")
}
