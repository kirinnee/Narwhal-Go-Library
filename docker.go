package narwhal_lib

type Docker struct {
	quiet bool
}

func (d Docker) ContainerIds() ([]string, []string) {

	psq := CreateCommand("docker", "ps", "-q")
	return d.containers(psq)

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
