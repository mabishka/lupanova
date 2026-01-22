package audit

import (
	"context"
	"encoding/json"

	"github.com/mabishka/lupanova/internal/model"
	"golang.org/x/sync/errgroup"
)

// Observer Аудит.
type Observer interface {
	GetName() string
	Send(context.Context, []byte) error
}

// AuditEvent события аудита.
type AuditEvent struct {
	observer map[string]Observer
}

// NewAuditEvent создание хранилища аудита.
func NewAuditEvent() *AuditEvent {
	return &AuditEvent{
		observer: make(map[string]Observer),
	}
}

// Register регистрация хранилища аудита.
func (p *AuditEvent) Register(o Observer) {
	if p.observer == nil {
		p.observer = make(map[string]Observer)
	}

	p.observer[o.GetName()] = o
}

// Send отправка аудита.
func (p *AuditEvent) Send(ctx context.Context, data *model.AuditData) error {

	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	wg, sendCtx := errgroup.WithContext(ctx)
	for _, v := range p.observer {
		wg.Go(func() error {
			return v.Send(sendCtx, resp)
		})
	}

	return wg.Wait()
}
