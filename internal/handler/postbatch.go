package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

// Эндпоинт /api/shorten/batch, принимающий в теле запроса множество URL для сокращения в формате json
func (p *StorageServer) HandlerPostBatch(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostBatch")
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
	var request []model.FullItem
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Error("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("request", zap.Int("count", len(request)))
	response, err := p.GetShortList(context.TODO(), request)
	if err != nil {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("response", zap.Int("count", len(response)))

	formatResponse := make([]model.ShortItem, 0, len(response))
	for _, v := range response {
		//v.Short = p.format(v.Short)

		formatResponse = append(formatResponse, model.ShortItem{Corr: v.Corr, Short: p.format(v.Short)})
		logger.Log().Info("item", zap.String("short", v.Short))
	}

	jsonResponse, err := json.Marshal(formatResponse)
	if err != nil {
		logger.Log().Error("Error marshal JSON response = ", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("response", zap.String("data", string(jsonResponse)))

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

}
