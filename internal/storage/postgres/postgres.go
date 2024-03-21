package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"testex/internal/config"
)

const (
	CommandTable          = "commands"
	ExecutedCommandsTable = "executed_commands"
	LogsTable             = "logs"
)

func New(cfg config.PostgresDatabase) (*sqlx.DB, error) {
	const fn = "storage.postgres.New"

	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, "disable"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS commands(
			id SERIAL PRIMARY KEY,
			alias varchar(128) UNIQUE,
		    script varchar(256)
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS executed_commands(
			id SERIAL PRIMARY KEY,
			command_id INT REFERENCES commands,
			PID INT NOT NULL,
			is_active BOOLEAN DEFAULT true
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS logs(
			id SERIAL PRIMARY KEY,
			executed_command_id INT REFERENCES executed_commands,
			message TEXT
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return db, nil
}
