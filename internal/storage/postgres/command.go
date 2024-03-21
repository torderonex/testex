package postgres

import (
	"fmt"
	"testex/internal/entities"

	"github.com/jmoiron/sqlx"
)

type CommandStorage struct {
	Db *sqlx.DB
}

func NewCommandStorage(db *sqlx.DB) *CommandStorage {
	return &CommandStorage{db}
}

func (s CommandStorage) SaveCommand(command entities.Command) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (alias ,script) VALUES ($1,$2) RETURNING id", CommandTable)
	row := s.Db.QueryRow(query, command.Alias, command.Script)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s CommandStorage) GetCommand(alias string) (entities.Command, error) {
	var c entities.Command
	query := fmt.Sprintf("SELECT * from %s WHERE alias=$1", CommandTable)
	err := s.Db.Get(&c, query, alias)
	return c, err
}

func (s CommandStorage) GetAllCommands() ([]entities.Command, error) {
	var c []entities.Command
	query := fmt.Sprintf("SELECT * from %s", CommandTable)
	err := s.Db.Select(&c, query)
	return c, err
}

func (s CommandStorage) SaveLog(log entities.Log) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (executed_command_id, message ) VALUES ($1,$2) RETURNING id", LogsTable)
	row := s.Db.QueryRow(query, log.ExecutedCommandId, log.Message)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s CommandStorage) SaveExecutedCommand(ec entities.ExecutedCommand) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (command_id, PID) VALUES ($1, $2) RETURNING id", ExecutedCommandsTable)
	row := s.Db.QueryRow(query, ec.CommandId, ec.PID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s CommandStorage) FinishCommand(commandID int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = $1 WHERE id = $2", ExecutedCommandsTable)
	_, err := s.Db.Exec(query, false, commandID)
	return err
}

func (s CommandStorage) GetExecutedCommandById(id int) (entities.ExecutedCommand, error) {
	var c entities.ExecutedCommand
	query := fmt.Sprintf("SELECT * from %s WHERE id = $1", ExecutedCommandsTable)
	err := s.Db.Get(&c, query, id)
	return c, err
}

func (s CommandStorage) GetLogsByExecutedCommand(executedCommandID int) ([]entities.Log, error) {
	var logs []entities.Log
	query := fmt.Sprintf("SELECT * from %s WHERE executed_command_id = $1", LogsTable)
	err := s.Db.Select(&logs, query, executedCommandID)
	return logs, err
}
