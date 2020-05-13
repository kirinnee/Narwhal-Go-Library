package narwhal_lib

type Docker struct {
	quiet bool
}

func (d Docker) ContainerIds() ([]string, string) {
	processes := make([]string, 0, 10)

	psq := CreateCommand("docker", "ps", "-q")
	err := psq.CustomRun(d.quiet, func(s string) {
		processes = append(processes, s)
	}, func(s string) {})
	if err != "" {
		return processes, err
	}
	return processes, ""

}
