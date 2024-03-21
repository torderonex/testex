package storage

import (
	"github.com/jmoiron/sqlx"
	"testex/internal/entities"
	"testex/internal/storage/postgres"
)

type Storage struct {
	CommandRepository
}

type CommandRepository interface {
	SaveCommand(command entities.Command) (int, error)
	GetCommand(alias string) (entities.Command, error)
	GetAllCommands() ([]entities.Command, error)
	SaveLog(entities.Log) (int, error)
	SaveExecutedCommand(entities.ExecutedCommand) (int, error)
	FinishCommand(commandID int) error
	GetLogsByExecutedCommand(executedCommandID int) ([]entities.Log, error)
	GetExecutedCommandById(id int) (entities.ExecutedCommand, error)
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		CommandRepository: postgres.NewCommandStorage(db),
	}
}
