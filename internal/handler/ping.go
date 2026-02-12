package handler

import (
	"net/http"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

// ConnServer сервер обработки пинга к БД.
type ConnServer struct {
	model.ConnLoader
}

// NewConn создание сервера обработки пинга к БД.
func NewConn(x model.ConnLoader) *ConnServer {
	return &ConnServer{ConnLoader: x}
}

// HandlerGetPing хендлер GET /ping, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
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
