package main

import (
	"net/http"

	"github.com/mabishka/lupanova/internal/handler"
)

const addr = "localhost:8080"

func main() {
	srv := handler.New(addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.HandlerPostFull)
	mux.HandleFunc("/{id}", srv.HandlerGetFull)

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
