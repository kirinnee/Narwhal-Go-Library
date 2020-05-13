package narwhal_lib

type Docker struct {
	quiet bool
}

func (d Docker) ContainerIds() ([]string, string) {
	containers := make([]string, 0, 10)

	psq := CreateCommand("docker", "ps", "-q")
	err := psq.CustomRun(d.quiet, func(s string) {
		containers = append(containers, s)
	}, func(s string) {})
	if err != "" {
		return containers, err
	}
	return containers, ""

}

func (d Docker) AllContainerIds() ([]string, string) {
	containers := make([]string, 0, 10)

	psq := CreateCommand("docker", "ps", "-aq")
	err := psq.CustomRun(d.quiet, func(s string) {
		containers = append(containers, s)
	}, func(s string) {})
	if err != "" {
		return containers, err
	}
	return containers, ""

}
