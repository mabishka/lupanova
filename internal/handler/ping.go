package handler

import (
	"net/http"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

type ConnServer struct {
	model.ConnLoader
}

func NewConn(x model.ConnLoader) *ConnServer {
	return &ConnServer{ConnLoader: x}
}

func (p *ConnServer) HandlerGetPing(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerGetPing")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := p.Ping(r.Context()); err != nil {
		logger.Log().Error("error ping", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeText)
	w.WriteHeader(http.StatusOK)
}
