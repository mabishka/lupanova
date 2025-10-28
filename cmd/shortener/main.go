package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mabishka/lupanova/internal/compress"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/repository/fileloader"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	server := handler.New(config.GetBaseAddress())
	loader := fileloader.New(config.GetFileName())
	server.Load(loader)
	router := chi.NewRouter()

	router.Post(`/`, logger.WithLogging(compress.WithCompress(server.HandlerPostFull)))
	router.Post(`/api/shorten`, logger.WithLogging(compress.WithCompress(server.HandlerPostFullJSON)))
	router.Get(`/{id}`, logger.WithLogging(compress.WithCompress(server.HandlerGetFull)))

	return http.ListenAndServe(config.GetServerAddress(), router)
}
