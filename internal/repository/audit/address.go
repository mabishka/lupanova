package audit

import (
	"bytes"
	"context"
	"net/http"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"go.uber.org/zap"
)

const observerAddressName = "url"

// AddressObserver отправка аудита по сети.
type AddressObserver struct {
	address string
}

// NewAddressObserver создание аудита для отправки по сети.
func NewAddressObserver(address string) *AddressObserver {
	return &AddressObserver{address: address}
}

// GetName имя аудита.
func (p *AddressObserver) GetName() string {
	return observerAddressName
}

// Send отправка аудита.
func (p *AddressObserver) Send(ctx context.Context, data []byte) error {

	r := bytes.NewReader(data)
	resp, err := http.Post(p.address, model.ContentTypeJSON, r)
	if err != nil {
		logger.Log().Error("send audit url", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	return nil
}
