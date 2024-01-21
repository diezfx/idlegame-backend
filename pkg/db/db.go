package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Querier interface {
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Config struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`

	MigrationsDir string `json:"migrationsDir"`
}

const connectionString = "postgres://%s:%s@%s:%d/%s"

type DB struct {
	cfg    Config
	client *sql.DB
}

func (db *DB) Close() error {
	return db.client.Close()
}

func (db *DB) Get(ctx context.Context, dest any, query string, args ...any) error {
	return sqlscan.Get(ctx, db.client, dest, query, args...)
}

func (db *DB) Select(ctx context.Context, dest any, query string, args ...any) error {
	return sqlscan.Select(ctx, db.client, dest, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.client.ExecContext(ctx, query, args...)
}

func New(cfg Config) (*DB, error) {
	connString := buildConnectionString(cfg)
	logger.Info(context.Background()).String("connectionString", connString).Msg("connectionString")
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	return &DB{cfg: cfg, client: conn}, nil
}

func buildConnectionString(cfg Config) string {
	return fmt.Sprintf(connectionString, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func (db *DB) Up(_ context.Context) error {
	driver, err := migratepgx.WithInstance(db.client, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", db.cfg.MigrationsDir), db.cfg.Database, driver)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

type txQuerier struct {
	getFunc    func(ctx context.Context, dest any, query string, args ...any) error
	selectFunc func(ctx context.Context, dest any, query string, args ...any) error
	execFunc   func(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (q *txQuerier) Get(ctx context.Context, dest any, query string, args ...any) error {
	return q.getFunc(ctx, dest, query, args...)
}

func (q *txQuerier) Select(ctx context.Context, dest any, query string, args ...any) error {
	return q.selectFunc(ctx, dest, query, args...)
}

func (q *txQuerier) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return q.execFunc(ctx, query, args...)
}

func (db *DB) WithTx(ctx context.Context, fn func(tx Querier) error) error {
	tx, err := db.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			logger.Error(ctx, err).Msg("failed transaction rollback")
		}
	}()

	// wrap tx in a querier anymous func
	wrappedTx := &txQuerier{
		getFunc: func(ctx context.Context, dest any, query string, args ...any) error {
			return db.Get(ctx, dest, query, args...)
		},
		selectFunc: func(ctx context.Context, dest any, query string, args ...any) error {
			return db.Select(ctx, dest, query, args...)
		},
		execFunc: func(ctx context.Context, query string, args ...any) (sql.Result, error) {
			return db.client.ExecContext(ctx, query, args...)
		},
	}

	if err := fn(wrappedTx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
