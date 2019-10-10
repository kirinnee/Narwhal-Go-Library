package narwhal_lib

import (
	"bufio"
	"fmt"
	"log"
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

func (command Command) Run(quiet bool) string {

	cmd := exec.Command(command.command, command.arg...)
	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		return err.Error()
	}

	// create scanner
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			if !quiet {
				s := fmt.Sprint(scanner.Text())
				log.Println(s)
			}
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
