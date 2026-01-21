package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

// Эндпоинт с методом GET и путём /{id},
// где id — идентификатор сокращённого URL (например, /EwHXdJfB).
// В случае успешной обработки запроса сервер возвращает ответ с кодом 307
// и оригинальным URL в HTTP-заголовке Location.
func (p *StorageServer) HandlerGetFull(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerGetFull")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	full, err := p.GetFull(context.TODO(), id)
	if err != nil {
		if errors.Is(err, model.ErrorDeleted) {
			logger.Log().Error("error getting full (is deleted)", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			return
		}
		logger.Log().Error("error getting full", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderLocation, full)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
