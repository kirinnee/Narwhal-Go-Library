package narwhal_lib

import (
	"context"
	"fmt"
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
		output []string
	}
)

func (c *TestCommandFactory) Create(command string, arg ...string) Executable {
	return &TestCommand{cmd: exec.Command(command, arg...), factory: c}
}

func (command *TestCommand) Write(b []byte) error {
	pipe, err := command.cmd.StdinPipe()
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

func (command *TestCommand) CustomRun(outEvent OutputEvent, errEvent OutputEvent) string {

	stdout, err := command.cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	stderr, err := command.cmd.StderrPipe()
	if err != nil {
		return err.Error()
	}

	pipe(stdout, outEvent)
	pipe(stderr, errEvent)

	err = command.cmd.Start()
	if err != nil {
		return err.Error()
	}
	err = command.cmd.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (command *TestCommand) Run() []string {
	errs := make([]string, 0, 10)
	others := command.CustomRun(func(s string) {
		command.factory.output = append(command.factory.output, s)
	}, func(s string) {
		errs = append(errs, s)
	})
	if others == "" {
		return errs
	}
	return append(errs, others)
}

func helpRun(command string, arg ...string) []string {
	ret := make([]string, 0, 10)
	f := CommandFactory{quiet: false}
	f.Create(command, arg...).CustomRun(func(s string) {
		ret = append(ret, s)
	}, func(s string) {
		ret = append(ret, s)
	})
	return ret
}

func helpRunPrint(command string) string {
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
