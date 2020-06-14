package narwhal_lib

import (
	"github.com/google/uuid"
	"gitlab.com/kiringo/narwhal_lib/command"
	"gitlab.com/kiringo/narwhal_lib/docker"
	"gitlab.com/kiringo/narwhal_lib/images"
	"log"
	"path/filepath"
)

type Narwhal struct {
	quiet  bool
	docker docker.Docker
	Cmd    command.Creator
}

func New(quiet bool) *Narwhal {
	return NewCustom(quiet, command.NewFactory(quiet))
}

func NewCustom(quiet bool, factory command.Creator) *Narwhal {
	return &Narwhal{
		quiet:  quiet,
		docker: docker.New(quiet, factory),
		Cmd:    factory,
	}
}

func (n *Narwhal) Print(s string) {
	if !n.quiet {
		log.Print(s)
	}
}

func (n *Narwhal) Load(volume string, tarPath string) []string {
	id, err := uuid.NewUUID()
	if err != nil {
		return []string{err.Error()}
	}
	var errors []string
	name := "narwhal-mount" + "abc" + id.String()
	c := Container{name, n.quiet, n.Cmd}

	n.Print("Creating container to connect to volume...")
	s := c.Start("alpine", volume, "/home/data", "dt")
	if len(s) != 0 {
		return s
	}

	n.Print("Copying to container")
	s = c.Copy(true, tarPath, "/home/")
	if len(s) != 0 {
		errors = append(errors, s...)
	} else {
		n.Print("Done copying!")
	}
	file := filepath.Base(tarPath)
	if len(errors) == 0 {
		n.Print("Renaming file")
		s = c.Exec("/home", "mv", file, "data.tar.gz")
		if len(s) != 0 {
			errors = append(errors, s...)
		} else {
			n.Print("Done renaming!")
		}
	}
	if len(errors) == 0 {
		n.Print("Unzipping volume")
		s = c.Exec("/home", "tar", "-xzf", "data.tar.gz")
		if len(s) != 0 {
			errors = append(errors, s...)
		} else {
			n.Print("Volume Unzipped!")
		}
	}

	n.Print("Killing Container...")
	s = c.Kill()
	if len(s) != 0 {
		errors = append(errors, s...)
		n.Print("Container failed to be killed!")
	} else {
		n.Print("Container Killed...")
	}

	n.Print("Removing Container...")
	s = c.Remove()
	if len(s) != 0 {
		errors = append(errors, s...)
		n.Print("Container failed to be removed!")
	} else {
		n.Print("Container remove...")
	}

	return errors
}

func (n *Narwhal) Save(volume string, tarName string, path string) []string {
	id, err := uuid.NewUUID()
	if err != nil {
		return []string{err.Error()}
	}
	var errors []string

	zipped := tarName + ".tar.gz"
	name := "narwhal-mount" + "abc" + id.String()
	c := Container{name, n.quiet, n.Cmd}

	n.Print("Creating container to connect to volume...")
	s := c.Start("alpine", volume, "/home/data", "dt")
	if len(s) != 0 {
		return s
	}
	n.Print("Container Created")

	n.Print("Zipping volume...")
	s = c.Exec("/home", "tar", "-czf", zipped, "data")
	if len(s) != 0 {
		errors = append(errors, s...)
	}
	n.Print("Volume Zipped!")

	if len(errors) == 0 {
		n.Print("Copying to host...")
		s = c.Copy(false, "/home/"+zipped, path)
		if len(s) != 0 {
			errors = append(errors, s...)
		}
		n.Print("Done copying!")
	}

	n.Print("Killing Container...")
	s = c.Kill()
	if len(s) != 0 {
		errors = append(errors, s...)
		n.Print("Container failed to be killed!")
	} else {
		n.Print("Container Killed...")
	}

	n.Print("Removing Container...")
	s = c.Remove()
	if len(s) != 0 {
		errors = append(errors, s...)
		n.Print("Container failed to be removed!")
	} else {
		n.Print("Container remove...")
	}

	return errors

}

func (n *Narwhal) KillAll() []string {

	containers, err := n.docker.ContainerIds()

	if len(err) != 0 {
		return err
	}

	if len(containers) == 0 {
		n.Print("No containers killed")
		return []string{}
	}

	n.Print("Killing all containers")

	containers = append([]string{"kill"}, containers...)

	kill := n.Cmd.Create("docker", containers...)
	errs := kill.Run()
	if len(errs) != 0 {
		return errs
	}
	return []string{}

}

func (n *Narwhal) RemoveAll() []string {
	containers, err := n.docker.AllContainerIds()

	if len(err) != 0 {
		return err
	}

	if len(containers) == 0 {
		n.Print("No containers removed")
		return []string{}
	}
	n.Print("Removing all containers")

	containers = append([]string{"rm"}, containers...)

	return n.Cmd.Create("docker", containers...).Run()
}

func (n *Narwhal) StopAll() []string {
	containers, err := n.docker.AllContainerIds()

	if len(err) != 0 {
		return err
	}

	if len(containers) == 0 {
		n.Print("No containers stopped")
		return []string{}
	}
	n.Print("Stopping all containers")

	containers = append([]string{"stop"}, containers...)

	return n.Cmd.Create("docker", containers...).Run()
}

func (n *Narwhal) Deploy(stack string, file string) []string {

	b, compose, ss, err := parse(file)

	if stack == "" {
		stack = ss
	}
	if stack == "" {
		return []string{"no stack name returned"}
	}

	if err != nil {
		return []string{err.Error()}
	}
	for k, v := range compose.Images {
		err := n.docker.Build(v.Context, v.File, k)
		if len(err) > 0 {
			return err
		}
	}
	deploy := n.Cmd.Create("docker", "stack", "deploy", "--prune", "--with-registry-auth", "--compose-file", "-", "--resolve-image", "changed", stack)
	err = deploy.Write(b)
	if err != nil {
		return []string{err.Error()}
	}
	return deploy.Run()
}

func (n *Narwhal) DeployAuto(stack string, file string, unsafe bool) []string {

	stackExist := n.Cmd.Create("docker", "stack", "ls").Run()

	if len(stackExist) > 0 {
		n.Print("Docker not in swarm mode... starting in swarm mode")
		err := n.Cmd.Create("docker", "swarm", "init").Run()
		if len(err) > 0 {
			return err
		}
	}

	errs := n.Deploy(stack, file)
	if len(errs) == 0 || !unsafe {
		return errs
	}
	for _, v := range errs {
		n.Print(v)
	}
	n.Print("Stack could not be deploy... automatically re-initialize swarm...")
	n.Cmd.Create("docker", "swarm", "leave", "--force").Run()
	n.Cmd.Create("docker", "swarm", "init").Run()
	return n.Deploy(stack, file)

}

func (n *Narwhal) Run(context, file, image, name string) []string {
	err := n.docker.Build(context, file, image)
	if len(err) > 0 {
		return err
	}
	return n.docker.Run(image, name)
}

func (n *Narwhal) Images(filter ...string) (images.Images, []string, []string) {
	label, dangle, ref, tq, remaining, err := parseFilters(filter...)
	if len(err) > 0 {
		return nil, nil, err
	}
	i, err := n.docker.Images(dangle, label, ref)
	if len(err) > 0 {
		return nil, nil, err
	}
	all, err := n.docker.Images("", "", []string{})
	queryParser, e := images.New(tq...)
	if e != nil {
		return nil, nil, []string{e.Error()}
	}
	queries, e := queryParser.Parse(all)
	if e != nil {
		return nil, nil, []string{e.Error()}
	}
	var x images.Images = i
	left := x.ProcessQuery(queries)
	return left, remaining, []string{}

}

func (n *Narwhal) StopStack(stack string, file string) []string {

	if file != "" {
		_, _, ss, err := parse(file)
		if err != nil {
			return []string{err.Error()}
		}
		if stack == "" {
			stack = ss
		}
	}

	if stack == "" {
		return []string{"no stack name returned"}
	}

	return n.Cmd.Create("docker", "stack", "rm", stack).Run()
}

func (n *Narwhal) MoveOut(ctx, file, image, from, to, command string) []string {

	err := n.docker.Build(ctx, file, image)
	if len(err) > 0 {
		return err
	}
	createDummy := n.Cmd.Create("docker", "create", "-ti", "--name", "narwhal-dummy", image, command)
	err = createDummy.Run()
	if len(err) > 0 {
		return err
	}
	cp := n.Cmd.Create("docker", "cp", "narwhal-dummy:"+from, to)
	err = cp.Run()
	if len(err) > 0 {
		return err
	}

	rm := n.Cmd.Create("docker", "rm", "-f", "narwhal-dummy")
	err = rm.Run()
	if len(err) > 0 {
		return err
	}
	rmi := n.Cmd.Create("docker", "rmi", image)
	return rmi.Run()
}

func (n *Narwhal) RemoveImage(filter ...string) []string {
	images, remain, err := n.Images(filter...)
	if len(err) > 0 {
		return err
	}

	if len(remain) > 0 {
		return append([]string{"unknown commands"}, remain...)
	}

	if len(images) == 0 {
		n.Print("no image was removed")
		return nil
	}

	args := []string{"rmi"}
	for _, v := range images {
		args = append(args, v.Name)
	}

	remove := n.Cmd.Create("docker", args...)
	return remove.Run()
}
