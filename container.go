package narwhal_lib

type Container struct {
	name  string
	quiet bool
}

func (c Container) Start(image string, mount string, mountTarget string, flags string) string {
	cmd := CreateCommand("docker", "run", "-"+flags, "--name", c.name, "-v", mount+":"+mountTarget, image)
	return cmd.Run(c.quiet)
}

func (c Container) Kill() string {
	cmd := CreateCommand("docker", "kill", c.name)
	return cmd.Run(c.quiet)
}

func (c Container) Remove() string {
	cmd := CreateCommand("docker", "rm", c.name)
	return cmd.Run(c.quiet)
}

func (c Container) Copy(into bool, from string, to string) string {
	if into {
		to = c.name + ":" + to
	} else {
		from = c.name + ":" + from
	}
	cmd := CreateCommand("docker", "cp", from, to)
	return cmd.Run(c.quiet)
}

func (c Container) Exec(workDir string, args ...string) string {
	a := []string{"exec", "-w", workDir, c.name}
	a = append(a, args...)
	cmd := CreateCommand("docker", a...)
	return cmd.Run(c.quiet)
}
