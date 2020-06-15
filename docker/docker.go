package docker

import (
	"gitlab.com/kiringo/narwhal_lib/command"
	"path/filepath"
)

type Docker struct {
	quiet bool
	cmd   command.Creator
}

func New(quiet bool, cmd command.Creator) Docker {
	return Docker{quiet: quiet, cmd: cmd}
}

func (d Docker) ContainerIds() ([]string, []string) {

	psq := d.cmd.Create("docker", "ps", "-q")
	return d.containers(psq)

}

func (d Docker) Build(context, file, image string) []string {
	file = filepath.Join(context, file)
	return d.cmd.Create("docker", "build", "--rm", "--tag", image, "--file", file, context).Run()
}

func (d Docker) Run(image, name, cmd string) []string {
	args := []string{
		"run",
		"--rm",
	}
	if name != "" {
		args = append(args, "--name")
		args = append(args, name)
	}
	args = append(args, image)
	if cmd != "" {
		args = append(args, cmd)
	}
	return d.cmd.Create("docker", args...).Run()
}

func (d Docker) containers(psq command.Executable) ([]string, []string) {
	containers := make([]string, 0, 10)
	errors := make([]string, 0, 10)
	err := psq.CustomRun(func(s string) {
		containers = append(containers, s)
	}, func(s string) {
		errors = append(errors, s)
	})
	if err != "" || len(errors) != 0 {
		return containers, append(errors, err)
	}
	return containers, errors
}

func (d Docker) AllContainerIds() ([]string, []string) {
	psq := d.cmd.Create("docker", "ps", "-aq")
	return d.containers(psq)
}
