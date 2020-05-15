package narwhal_lib

import (
	"fmt"
	"path/filepath"
)

type Docker struct {
	quiet bool
	cmd   CommandCreator
}

func (d Docker) ContainerIds() ([]string, []string) {

	psq := d.cmd.Create("docker", "ps", "-q")
	return d.containers(psq)

}

func (d Docker) Build(context, file, image string) []string {
	file = filepath.Join(context, file)
	fmt.Println("HELOOOOOOOO", d.cmd)
	return d.cmd.Create("docker", "build", "--tag", image, "--file", file, context).Run()
}

func (d Docker) Run(image string) []string {
	return d.cmd.Create("docker", "run", "--rm", image).Run()
}

func (d Docker) containers(psq Executable) ([]string, []string) {
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
