package entities

import (
	"os/exec"
	"time"
)

type Command struct {
	Id     int    `json:"id"`
	Alias  string `json:"alias"`
	Script string `json:"script"`
}

type ExecutedCommand struct {
	Id        int
	CommandId int `db:"command_id"`
	PID       int
	IsActive  bool `db:"is_active"`
}

type Log struct {
	Id                int       `db:"id" json:"id"`
	ExecutedCommandId int       `db:"executed_command_id" json:"executed_command_id"`
	Message           string    `db:"message" json:"message"`
	Date              time.Time `db:"date" json:"date"`
}

type CommandInfo struct {
	Id       int
	Command  string
	Process  *exec.Cmd
	Finished bool
}

type CommandDto struct {
	Alias  string `json:"alias"`
	Script string `json:"script"`
}

type CommandIDResponse struct {
	Id int `json:"id"`
}

type ExecuteCommandResponse struct {
	Output string `json:"output"`
}
