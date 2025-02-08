package run

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Error struct {
	baseError error
	stdErr string
}

func errorFrom(baseError error, stdErr string) Error {
	return Error{
		baseError,
		stdErr,
	}
}

func (self *Error) Error() string {
	return fmt.Sprintf(
		"Command run error (%s)\n" +
			"CMD ERR OUTPUT:\n%s",
		self.baseError.Error(),
		self.stdErr,
	)
}

func (self *Error) ErrorOutput() string {
	return self.stdErr
}

type cmdRun struct {
	name string
	args []string
	timeoutSeconds *int
}

type cmdRunOption func(*cmdRun)

func NewCmd(name string, ops ...cmdRunOption) *cmdRun {
	c := &cmdRun{
		name: name,
		args: make([]string, 0),
		timeoutSeconds: nil,
	}

	for _, op := range ops {
		op(c)
	}

	return c
}

func WithArgs(args ...string) cmdRunOption {
	return func(c *cmdRun) {
		c.args = append(c.args, args...)
	}
}

func WithTimeout(seconds int) cmdRunOption {
	return func(c *cmdRun) {
		c.timeoutSeconds = &seconds
	}
}

func (self *cmdRun) Exec() error {
	var cmd *exec.Cmd
	if self.timeoutSeconds == nil {
		cmd = exec.Command(self.name, self.args...)
	} else {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(*self.timeoutSeconds) * time.Second,
		)
		defer cancel()
		cmd = exec.CommandContext(ctx, self.name, self.args...)
	}

	stdErr := &strings.Builder{}
	cmd.Stderr = stdErr

	err := cmd.Run()
	if err != nil {
		cmdError := errorFrom(err, stdErr.String())
		return &cmdError
	}
	return nil
}
