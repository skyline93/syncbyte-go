package dest

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
)

type PostgreSQL struct {
	Opts *Options
}

func NewPostgreSQL(opts *Options) *PostgreSQL {
	return &PostgreSQL{Opts: opts}
}

func (s *PostgreSQL) AdminUri(dbName string) string {
	return fmt.Sprintf(
		"postgresql://postgres:%s@127.0.0.1:%d/%s",
		s.Opts.Password, s.Opts.Port, dbName,
	)
}

func (p *PostgreSQL) Build(destFile string) error {
	if err := p.createEmptyDB(); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	if err := p.initDB(); err != nil {
		return err
	}

	if err := p.restore(destFile); err != nil {
		return err
	}

	if err := p.activeDB(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) createEmptyDB() error {
	c := fmt.Sprintf(
		"docker run -d -p %d:5432 --restart=always -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=postgres postgres:%s",
		p.Opts.Port, p.Opts.Password, p.Opts.Version,
	)
	_, err := utils.Exec("/bin/sh", "-c", c)
	return err
}

func (p *PostgreSQL) initDB() error {
	conn, err := pgx.Connect(context.Background(), p.AdminUri("postgres"))

	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s;", p.Opts.DBName))

	if _, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", p.Opts.User, p.Opts.Password)); err != nil {
		return err
	}

	if _, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s OWNER %s;", p.Opts.DBName, p.Opts.User)); err != nil {
		return err
	}

	if _, err := conn.Exec(context.Background(), fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", p.Opts.DBName, p.Opts.User)); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) restore(destFile string) error {
	c := fmt.Sprintf("pg_restore --no-owner --dbname=%s %s", p.AdminUri(p.Opts.DBName), destFile)
	_, err := utils.Exec("/bin/sh", "-c", c)
	return err
}

func (p *PostgreSQL) activeDB() error {
	conn, err := pgx.Connect(context.Background(), p.AdminUri(p.Opts.DBName))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), fmt.Sprintf("GRANT ALL PRIVILEGES ON all tables in schema public TO %s;", p.Opts.User)); err != nil {
		return err
	}

	return nil
}
