package lib

import (
	"bytes"
	"os/exec"

	"github.com/rs/zerolog"
)

type commandExecutor struct {
	log *zerolog.Logger
}

type CommandExecutor interface {
	Command(name string, args ...string) ([]byte, error)
}

// A thin wrapper around exec.Command
func NewCommandExecutor(l *zerolog.Logger) CommandExecutor {
	return &commandExecutor{
		log: l,
	}
}

// Executes a command and returns the output as a buffer of bytes
func (c *commandExecutor) Command(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	c.log.Info().Msgf("Executing command: %s %s", cmd, args)
	if err := cmd.Run(); err != nil {
		c.log.Error().Err(err).Msgf("Failed to execute command: %s %s", cmd, args)
		return nil, err
	}

	c.log.Info().Msgf("Command executed successfully: %s %s", cmd, args)

	return out.Bytes(), nil
}
