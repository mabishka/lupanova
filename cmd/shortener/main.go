package main

import (
	"context"
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
	run(context.Background())
}

func run(ctx context.Context) {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	server := handler.New(config.GetBaseAddress())

	var loader model.StorageLoader
	var conn model.ConnLoader
	if config.GetConnAddress() != "" {
		loader = connloader.New(config.GetConnAddress())
		conn, _ = loader.(model.ConnLoader)
		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Error("conn not loaded", zap.Error(err))
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

	router.Post(`/`, logger.WithLogging(compress.WithCompress(server.HandlerPostFull)))
	router.Post(`/api/shorten`, logger.WithLogging(compress.WithCompress(server.HandlerPostFullJSON)))
	router.Get(`/{id}`, logger.WithLogging(compress.WithCompress(server.HandlerGetFull)))
	router.Get(`/ping`, logger.WithLogging(compress.WithCompress(connServer.HandlerGetPing)))

	go func() {
		if err := http.ListenAndServe(config.GetServerAddress(), router); err != nil {
			panic(err)
		}
	}()
	<-ctx.Done()
}
