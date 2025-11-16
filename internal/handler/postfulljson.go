package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса JSON URL как application/json
// и возвращает ответ с кодом 201 и сокращённым JSON URL как application/json.
func (p *StorageServer) HandlerPostFullJSON(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostFullJSON")

	if r.Method != http.MethodPost {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("error contect type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	var request model.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Error("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	full := strings.TrimSpace(string(request.Full))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := p.GetShort(context.TODO(), full)
	if err != nil {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := model.Response{
		Short: p.format(short),
	}

	enc, err := json.Marshal(response)
	if err != nil {
		logger.Log().Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(enc)
}
