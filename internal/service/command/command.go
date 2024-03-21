package command

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"testex/internal/config"
	"testex/internal/entities"
	"testex/internal/storage"
)

type Service struct {
	Storage *storage.Storage
	Logger  *slog.Logger
	Config  config.Config
	mutex   sync.Mutex
}

func NewService(storage *storage.Storage, logger *slog.Logger, cfg config.Config) *Service {
	return &Service{
		Storage: storage,
		Logger:  logger,
		Config:  cfg,
	}
}

func (c *Service) Create(alias string, script string) (int, error) {
	return c.Storage.SaveCommand(entities.Command{Alias: alias, Script: script})
}

func (c *Service) GetAll() ([]entities.Command, error) {
	return c.Storage.GetAllCommands()
}

func (c *Service) GetOne(alias string) (entities.Command, error) {
	return c.Storage.GetCommand(alias)
}

func (c *Service) Execute(alias string) (int, error) {
	ctx := context.Background()
	var (
		name = "bash"
		arg  = "-c"
	)

	if c.Config.Os == "win" {
		name = "cmd"
		arg = "/C"
	}

	command, err := c.GetOne(alias)
	if err != nil {
		return -1, err
	}

	cmd := exec.CommandContext(ctx, name, arg, command.Script)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	err = cmd.Start()
	if err != nil {
		return -1, err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	id, err := c.Storage.SaveExecutedCommand(entities.ExecutedCommand{CommandId: command.Id, PID: cmd.Process.Pid})
	if err != nil {
		return -1, err
	}

	go func() {
		stdoutScanner := bufio.NewScanner(stdout)
		stderrScanner := bufio.NewScanner(stderr)

		go func() {
			for stdoutScanner.Scan() {
				msg := fmt.Sprintf("[%d - STDOUT] %s\n", id, stdoutScanner.Text())
				c.Logger.Info(stdoutScanner.Text())
				_, err = c.Storage.SaveLog(entities.Log{Message: msg, ExecutedCommandId: id})
				if err != nil {
					c.Logger.Error(err.Error())
				}
			}
		}()

		go func() {
			for stderrScanner.Scan() {
				msg := fmt.Sprintf("[%d - STDERR] %s\n", id, stderrScanner.Text())
				c.Logger.Error(msg)
				c.Storage.SaveLog(entities.Log{Message: msg, ExecutedCommandId: id})
			}
		}()

		err := cmd.Wait()
		if err != nil {
			log.Printf("Command failed: %v", err)
		}

		c.mutex.Lock()
		defer c.mutex.Unlock()
		_ = c.Storage.FinishCommand(id)
	}()

	return id, err
}

func (c *Service) StopCommand(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	cmd, err := c.Storage.GetExecutedCommandById(id)
	if err != nil {
		return err
	}

	if !cmd.IsActive {
		return fmt.Errorf("command is not active")
	}

	cmdInfo, err := os.FindProcess(cmd.PID)
	if err != nil {
		return err
	}

	err = cmdInfo.Kill()
	if err != nil {
		return err
	}

	err = c.Storage.FinishCommand(id)

	return err
}

func (c *Service) GetLogs(executedCommandId int) ([]entities.Log, error) {
	return c.Storage.GetLogsByExecutedCommand(executedCommandId)
}
