package service

import (
	"log/slog"
	"testex/internal/config"
	"testex/internal/entities"
	"testex/internal/service/command"
	"testex/internal/storage"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	Command
}

func New(s *storage.Storage, logger *slog.Logger, cfg config.Config) *Service {
	return &Service{
		Command: command.NewService(s, logger, cfg),
	}
}

type Command interface {
	Execute(alias string) (int, error)
	Create(alias string, script string) (int, error)
	GetAll() ([]entities.Command, error)
	GetOne(alias string) (entities.Command, error)
	StopCommand(id int) error
	GetLogs(executedCommandId int) ([]entities.Log, error)
}
