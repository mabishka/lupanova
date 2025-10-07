package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mabishka/lupanova/internal/handler"
)

const addr = "localhost:8080"

func main() {
	server := handler.New(addr)
	router := chi.NewRouter()
	router.Post(`/`, server.HandlerPostFull)
	router.Get(`/{id}`, server.HandlerGetFull)

	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
