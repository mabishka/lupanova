package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/pkg/utils"
	"go.uber.org/zap"
)

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса строку URL как text/plain
// и возвращает ответ с кодом 201 и сокращённым URL как text/plain.
func (p *StorageServer) HandlerPostFull(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostFull")
	if r.Method != http.MethodPost {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeText {
		logger.Log().Error("error content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("error getting request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	full := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, shorterr := p.GetShort(context.TODO(), full, getUser(r))
	if shorterr != nil && !errors.Is(shorterr, utils.ErrConflict) {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeText)
	if shorterr != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(p.format(short)))
}
