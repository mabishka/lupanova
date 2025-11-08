package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mabishka/lupanova/internal/compress"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
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

	loader := fileloader.New(config.GetFileName())
	if err := server.Load(loader); err != nil {
		logger.Log().Info("memory storage usage")
	} else {
		logger.Log().Info("file storage usage")
	}

	conn := connloader.New(config.GetConnAddress())
	if err := conn.Load(context.TODO()); err == nil {
		logger.Log().Error("conn not loaded")
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
