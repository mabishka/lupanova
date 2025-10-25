package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
)

func main() {

	config := config.New()
	server := handler.New(config.GetBaseAddress())

	router := chi.NewRouter()
	router.Post(`/`, server.HandlerPostFull)
	router.Get(`/{id}`, server.HandlerGetFull)

	fmt.Println(config)

	if err := http.ListenAndServe(config.GetServerAddress(), router); err != nil {
		panic(err)
	}
}
