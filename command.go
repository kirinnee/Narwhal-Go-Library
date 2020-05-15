package narwhal_lib

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type (
	Command struct {
		cmd   *exec.Cmd
		quiet bool
	}

	CommandFactory struct {
		quiet bool
	}
)

func (c CommandFactory) Create(command string, arg ...string) Executable {
	return &Command{quiet: c.quiet, cmd: exec.Command(command, arg...)}

}

func (command *Command) Write(b []byte) error {
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

func pipe(pipe io.ReadCloser, event OutputEvent) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			s := fmt.Sprint(scanner.Text())
			event(s)
		}
	}()

}

func (command *Command) CustomRun(outEvent OutputEvent, errEvent OutputEvent) string {

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

func (command *Command) Run() []string {
	errs := make([]string, 0, 10)
	others := command.CustomRun(func(s string) {
		if !command.quiet {
			fmt.Println(s)
		}
	}, func(s string) {
		errs = append(errs, s)
	})
	if others == "" {
		return errs
	}
	return append(errs, others)
}
