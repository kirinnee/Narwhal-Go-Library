package test_helper

import (
	"context"
	"fmt"
	"gitlab.com/kiringo/narwhal_lib/command"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"os/exec"
	"strings"
)

type (
	TestCommand struct {
		cmd     *exec.Cmd
		factory *TestCommandFactory
	}

	TestCommandFactory struct {
		Output []string
	}
)

func (c *TestCommandFactory) Create(cmd string, arg ...string) command.Executable {
	return &TestCommand{cmd: exec.Command(cmd, arg...), factory: c}
}

func (c *TestCommand) Write(b []byte) error {
	pipe, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, err = pipe.Write(b)
	if err != nil {
		return err
	}
	err = pipe.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *TestCommand) CustomRun(outEvent command.OutputEvent, errEvent command.OutputEvent) string {

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return err.Error()
	}

	command.Pipe(stdout, outEvent)
	command.Pipe(stderr, errEvent)

	err = c.cmd.Start()
	if err != nil {
		return err.Error()
	}
	err = c.cmd.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (c *TestCommand) Run() []string {
	errs := make([]string, 0, 10)
	others := c.CustomRun(func(s string) {
		c.factory.Output = append(c.factory.Output, s)
	}, func(s string) {
		errs = append(errs, s)
	})
	if others == "" {
		return errs
	}
	return append(errs, others)
}

func HelpRun(c string, arg ...string) []string {
	ret := make([]string, 0, 10)
	f := TestCommandFactory{[]string{}}
	f.Create(c, arg...).CustomRun(func(s string) {
		ret = append(ret, s)
	}, func(s string) {
		ret = append(ret, s)
	})
	return ret
}

func HelpRunQ(command string) string {
	ctx := context.Background()
	runner, _ := interp.New(interp.StdIO(nil, nil, nil))

	f, _ := syntax.NewParser().Parse(strings.NewReader(command), "")

	err := runner.Run(ctx, f)
	if err != nil {
		return err.Error()
	}
	return ""
}

func HelpRunPrint(command string) string {
	ctx := context.Background()
	runner, _ := interp.New(interp.StdIO(nil, LogWriter{}, LogWriter{}))

	f, _ := syntax.NewParser().Parse(strings.NewReader(command), "")

	err := runner.Run(ctx, f)
	if err != nil {
		return err.Error()
	}
	return ""
}

type LogWriter struct {
}

func (LogWriter) Write(p []byte) (n int, err error) {

	fmt.Println(string(p))
	return len(p), nil

}
