package narwhal_lib

import "path/filepath"

type Docker struct {
	quiet bool
}

func (d Docker) ContainerIds() ([]string, []string) {

	psq := CreateCommand("docker", "ps", "-q")
	return d.containers(psq)

}

func (d Docker) Build(context, file, image string) []string {
	file = filepath.Join(context, file)
	return CreateCommand("docker", "build", "--tag", image, "--file", file, context).Run(d.quiet)
}

func (d Docker) Run(image string) []string {
	return CreateCommand("docker", "run", "--rm", image).Run(d.quiet)
}

func (d Docker) containers(psq Command) ([]string, []string) {
	containers := make([]string, 0, 10)
	errors := make([]string, 0, 10)
	err := psq.CustomRun(d.quiet, func(s string) {
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
	psq := CreateCommand("docker", "ps", "-aq")
	return d.containers(psq)
}
