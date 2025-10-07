package handler

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mabishka/lupanova/internal/service"
)

type StorageServer struct {
	service.Storage
	u url.URL
}

func New(addr string) *StorageServer {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	if host == "" {
		host = "localhost"
	}
	u := url.URL{
		Host:   net.JoinHostPort(host, port),
		Scheme: "http",
	}
	return &StorageServer{Storage: service.New(), u: u}
}

func (p *StorageServer) format(path string) string {
	p.u.Path = path
	return p.u.String()
}

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса строку URL как text/plain
// и возвращает ответ с кодом 201 и сокращённым URL как text/plain.
func (p *StorageServer) HandlerPostFull(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	full := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(full); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := p.GetShort(full)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(p.format(short)))
}

// Эндпоинт с методом GET и путём /{id},
// где id — идентификатор сокращённого URL (например, /EwHXdJfB).
// В случае успешной обработки запроса сервер возвращает ответ с кодом 307
// и оригинальным URL в HTTP-заголовке Location.
func (p *StorageServer) HandlerGetFull(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	full, err := p.GetFull(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", full)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
