package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
)

func main() {

	config := config.New()

	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	server := handler.New(config.GetBaseAddress())
	router := chi.NewRouter()

	router.Post(`/`, logger.WithLogging(server.HandlerPostFull))
	router.Get(`/{id}`, logger.WithLogging(server.HandlerGetFull))

	if err := http.ListenAndServe(config.GetServerAddress(), router); err != nil {
		panic(err)
	}
}
