package connloader

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/repository/db"
	"github.com/mabishka/lupanova/pkg/utils"
	"go.uber.org/zap"
)

type ConnLoader struct {
	conn *sql.DB
	addr string
}

func New(addr string) *ConnLoader {
	return &ConnLoader{addr: addr}

}

func (p *ConnLoader) Create(ctx context.Context) error {

	db, err := sql.Open("pgx", p.addr)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(2) // Установить максимальное количество открытых соединений к базе данных
	db.SetMaxIdleConns(2) // Установить максимальное количество неактивных соединений в пуле

	p.conn = db
	return nil
}

func (p *ConnLoader) Ping(ctx context.Context) error {

	if p.conn == nil {
		if err := p.Create(ctx); err != nil {
			logger.Log().Error("error", zap.Error(err))
			return err
		}
	}
	return p.conn.PingContext(ctx)
}

func (p *ConnLoader) Load(ctx context.Context) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	if err := db.Create(ctx, p.conn); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	return db.LoadList(ctx, p.conn)
}

func (p *ConnLoader) GetShortList(ctx context.Context, fullList []model.FullItem, user string) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	shortList := make(map[string]string)

	for _, full := range fullList {

		short, err := db.GetShort(ctx, tx, full.Full, user)
		if err != nil {
			if !errors.Is(err, utils.ErrConflict) {
				logger.Log().Error("error", zap.Error(err))
				tx.Rollback()
				return nil, err
			}
		}
		shortList[full.Full] = short
	}
	return shortList, tx.Commit()
}

func (p *ConnLoader) GetShort(ctx context.Context, full string, user string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	short, err := db.GetShort(ctx, tx, full, user)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		tx.Rollback()
		return short, err
	}

	return short, tx.Commit()
}

func (p *ConnLoader) GetFull(ctx context.Context, short string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	return db.GetFull(ctx, p.conn, short)
}

func (p *ConnLoader) GetUserList(ctx context.Context, user string) ([]model.StoreItem, error) {
	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	return db.GetUser(ctx, p.conn, user)
}
