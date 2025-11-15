package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/pkg/rand"
)

type Connector interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func Create(ctx context.Context, conn Connector) error {

	_, err := conn.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS t_data (
    		id SERIAL PRIMARY KEY,
    		s_full VARCHAR(1000) NOT NULL,
    		s_short VARCHAR(100) NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_data_full ON t_data(s_full);
		CREATE INDEX IF NOT EXISTS idx_data_short ON t_data(s_short); `)

	return err
}

func LoadList(ctx context.Context, conn Connector) (map[string]string, error) {

	rows, err := conn.QueryContext(ctx, "select s_full, s_short from t_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var full, short *string

	list := make(map[string]string)
	for rows.Next() {
		if err := rows.Scan(&full, &short); err != nil {
			return nil, err
		}
		if full != nil && short != nil {
			list[*short] = *full
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func store(ctx context.Context, conn Connector, full, short string) error {
	_, err := conn.ExecContext(ctx, "insert into t_data(s_full, s_short) values($1, $2)", full, short)
	return err
}

func GetShort(ctx context.Context, conn Connector, full string) (string, error) {
	if short, err := getShort(ctx, conn, full); err == nil {
		return short, nil
	}

	short, err := rand.CreateShort(config.ShortLen)
	if err != nil {
		return "", err
	}

	if err := store(ctx, conn, full, short); err != nil {
		return "", err
	}
	return short, nil
}

func getShort(ctx context.Context, conn Connector, full string) (string, error) {
	rows, err := conn.QueryContext(ctx, "select s_short from t_data where s_full = $1", full)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("full name for %s not found", full)
	}

	var short *string
	if err := rows.Scan(&short); err != nil {
		return "", err
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	if short == nil {
		return "", fmt.Errorf("full name for %s is empty", full)
	}

	return *short, nil
}

func GetFull(ctx context.Context, conn Connector, short string) (string, error) {
	rows, err := conn.QueryContext(ctx, "select s_full from t_data where s_short = $1", short)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("full name for %s not found", short)
	}

	var full *string
	if err := rows.Scan(&full); err != nil {
		return "", err
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	if full == nil {
		return "", fmt.Errorf("full name for %s is empty", short)
	}

	return *full, nil
}
