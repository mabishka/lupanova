package connloader

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/repository/db"
)

type ConnLoader struct {
	conn *sql.DB
	addr string
}

func New(addr string) *ConnLoader {
	return &ConnLoader{addr: addr}

}

func (p *ConnLoader) Create(ctx context.Context) error {

	db, err := sql.Open("postgres", p.addr)
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	p.conn = db
	return nil
}

func (p *ConnLoader) Ping(ctx context.Context) error {

	if p.conn == nil {
		if err := p.Create(ctx); err != nil {
			return err
		}
	}
	return p.conn.PingContext(ctx)
}

func (p *ConnLoader) Load(ctx context.Context) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	if err := db.Create(ctx, p.conn); err != nil {
		return nil, err
	}

	return db.LoadList(ctx, p.conn)
}

func (p *ConnLoader) GetShortList(ctx context.Context, fullList []model.FullItem) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	shortList := make(map[string]string)

	for _, full := range fullList {

		short, err := db.GetShort(ctx, tx, full.Full)
		if err != nil {
			return nil, err
		}
		shortList[full.Full] = short
	}
	return shortList, tx.Commit()
}

func (p *ConnLoader) GetShort(ctx context.Context, full string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		return "", err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	short, err := db.GetShort(ctx, tx, full)
	if err != nil {
		return "", err
	}

	return short, tx.Commit()
}

func (p *ConnLoader) GetFull(ctx context.Context, short string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		return "", err
	}

	return db.GetFull(ctx, p.conn, short)
}
