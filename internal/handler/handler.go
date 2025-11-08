package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/service"
)

type StorageServer struct {
	model.Storage
	u *url.URL
}

func New(address string) *StorageServer {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	return &StorageServer{Storage: service.New(), u: u}
}

func (p *StorageServer) SetLoader(loader model.Storage) {
	p.Storage = loader
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
		logger.Log().Debug("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeText {
		logger.Log().Debug("error content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Debug("error getting request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	full := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Debug("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := p.GetShort(context.TODO(), full)
	if err != nil {
		logger.Log().Debug("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeText)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(p.format(short)))
}

// Эндпоинт с методом GET и путём /{id},
// где id — идентификатор сокращённого URL (например, /EwHXdJfB).
// В случае успешной обработки запроса сервер возвращает ответ с кодом 307
// и оригинальным URL в HTTP-заголовке Location.
func (p *StorageServer) HandlerGetFull(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		logger.Log().Debug("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	full, err := p.GetFull(context.TODO(), id)
	if err != nil {
		logger.Log().Debug("error getting full", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderLocation, full)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса JSON URL как application/json
// и возвращает ответ с кодом 201 и сокращённым JSON URL как application/json.
func (p *StorageServer) HandlerPostFullJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log().Debug("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Debug("error contect type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	var request model.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Debug("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	full := strings.TrimSpace(string(request.Full))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Debug("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := p.GetShort(context.TODO(), full)
	if err != nil {
		logger.Log().Debug("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := model.Response{
		Short: p.format(short),
	}

	enc, err := json.Marshal(response)
	if err != nil {
		logger.Log().Debug("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(enc)
}

type ConnServer struct {
	model.ConnLoader
}

func NewConn(x model.ConnLoader) *ConnServer {
	return &ConnServer{ConnLoader: x}
}

func (p *ConnServer) HandlerGetPing(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		logger.Log().Debug("error method")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := p.Ping(context.TODO()); err != nil {
		logger.Log().Debug("error ping", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
}
