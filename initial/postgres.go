package initial

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"GinChat/config"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() error {
	cfg := config.Global.Postgres
	if err := ensureDatabase(cfg); err != nil {
		return err
	}

	db, err := openPostgres(cfg, cfg.DBName)
	if err != nil {
		return fmt.Errorf("connect to postgres database %q: %w", cfg.DBName, err)
	}

	config.Global.DB = db
	return nil
}

func openPostgres(cfg config.Postgres, dbName string) (*gorm.DB, error) {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:     "/" + dbName,
		RawQuery: "sslmode=disable",
	}
	return gorm.Open(postgres.Open(dsn.String()))
}

func ensureDatabase(cfg config.Postgres) error {
	serverDB, err := openPostgres(cfg, "postgres")
	if err != nil {
		return fmt.Errorf("connect to postgres maintenance database: %w", err)
	}
	sqlDB, err := serverDB.DB()
	if err != nil {
		return fmt.Errorf("get postgres maintenance connection: %w", err)
	}
	defer sqlDB.Close()

	var exists bool
	if err := serverDB.Raw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)",
		cfg.DBName).Row().Scan(&exists); err != nil {
		return fmt.Errorf("check postgres database %q: %w", cfg.DBName, err)
	}
	if exists {
		return nil
	}

	if err := serverDB.Exec("CREATE DATABASE " + pgx.Identifier{cfg.DBName}.Sanitize()).Error; err != nil {
		return fmt.Errorf("create postgres database %q: %w", cfg.DBName, err)
	}
	return nil
}
