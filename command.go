package narwhal_lib

import (
	"bufio"
	"fmt"
	"os/exec"
)

type Command struct {
	command string
	arg     []string
}

type OutputEvent = func(string)

func CreateCommand(command string, arg ...string) Command {
	return Command{command, arg}
}

func (command Command) CustomRun(quiet bool, outEvent OutputEvent, errEvent OutputEvent) string {
	cmd := exec.Command(command.command, command.arg...)

	// Error Scan
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err.Error()
	}

	// Normal Scan
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}

	// create scanners for both pipes
	errScanner := bufio.NewScanner(errReader)
	outScanner := bufio.NewScanner(outReader)

	go func() {
		for errScanner.Scan() {
			s := fmt.Sprint(errScanner.Text())
			errEvent(s)
		}
	}()
	go func() {
		for outScanner.Scan() {
			s := fmt.Sprint(outScanner.Text())
			outEvent(s)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err.Error()
	}

	err = cmd.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (command Command) Run(quiet bool) []string {
	errs := make([]string, 0, 10)
	others := command.CustomRun(quiet, func(s string) {
		if !quiet {
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
