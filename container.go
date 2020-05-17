package narwhal_lib

import "gitlab.com/kiringo/narwhal_lib/command"

type Container struct {
	name  string
	quiet bool
	cmd   command.Creator
}

func (c Container) Start(image string, mount string, mountTarget string, flags string) []string {
	cmd := c.cmd.Create("docker", "run", "-"+flags, "--name", c.name, "-v", mount+":"+mountTarget, image)
	return cmd.Run()
}

func (c Container) Kill() []string {
	cmd := c.cmd.Create("docker", "kill", c.name)
	return cmd.Run()
}

func (c Container) Remove() []string {
	cmd := c.cmd.Create("docker", "rm", c.name)
	return cmd.Run()
}

func (c Container) Copy(into bool, from string, to string) []string {
	if into {
		to = c.name + ":" + to
	} else {
		from = c.name + ":" + from
	}
	cmd := c.cmd.Create("docker", "cp", from, to)
	return cmd.Run()
}

func (c Container) Exec(workDir string, args ...string) []string {
	a := []string{"exec", "-w", workDir, c.name}
	a = append(a, args...)
	cmd := c.cmd.Create("docker", a...)
	return cmd.Run()
}
