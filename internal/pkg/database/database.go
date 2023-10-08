package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/anhnmt/golang-grpc-base-project/ent"
	"github.com/anhnmt/golang-grpc-base-project/ent/migrate"
	"github.com/anhnmt/golang-grpc-base-project/internal/pkg/config"
)

func New(ctx context.Context) (*ent.Client, error) {
	if !config.DatabaseEnabled() {
		return nil, nil
	}

	maxOpenConns := config.DatabaseMaxOpenConns()
	if maxOpenConns == 0 {
		maxOpenConns = 15
	}

	maxIdleConns := config.DatabaseMaxIdleConns()
	if maxIdleConns == 0 {
		maxIdleConns = 2
	}

	maxLifetime := config.DatabaseMaxLifetime()
	if maxLifetime == 0 {
		maxLifetime = 5
	}

	host := config.DatabaseHost()
	port := config.DatabasePort()
	pgbouncer := config.DatabasePgbouncer()

	log.Info().
		Str("host", host).
		Int("port", port).
		Bool("pgbouncer", pgbouncer).
		Int("maxOpenConns", maxOpenConns).
		Int("maxIdleConns", maxIdleConns).
		Int("maxLifetime", maxLifetime).
		Msg("Connecting to DB")

	dsn := &url.URL{
		Scheme: dialect.Postgres,
		User:   url.UserPassword(config.DatabaseUser(), config.DatabasePassword()),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   config.DatabaseName(),
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	cfg, err := pgx.ParseConfig(dsn.String())
	if err != nil {
		return nil, err
	}

	if config.DatabasePgbouncer() {
		cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

	db := stdlib.OpenDB(*cfg)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(maxLifetime) * time.Minute)

	newCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err = db.PingContext(newCtx)
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	// Create an ent.Driver from `db`.
	client := ent.NewClient(ent.Driver(drv))
	if config.DatabaseDebug() {
		client = client.Debug()
	}

	// Run the auto migration tool.
	if config.DatabaseMigration() {
		if err = client.Schema.Create(
			ctx,
			migrate.WithForeignKeys(false), // Disable foreign keys.
		); err != nil {
			log.Err(err).Msg("Failed creating schema resources")
			return nil, err
		}

		log.Info().Msg("Migrate to DB successfully.")
	}

	log.Info().Msg("Connecting to DB successfully.")
	return client, nil
}
