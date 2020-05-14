package narwhal_lib

import (
	"github.com/google/uuid"
	"log"
	"path/filepath"
)

type Narwhal struct {
	quiet  bool
	docker Docker
}

func New(quiet bool) *Narwhal {
	return &Narwhal{
		quiet:  quiet,
		docker: Docker{quiet: quiet},
	}
}

func (n Narwhal) Print(s string) {
	if !n.quiet {
		log.Print(s)
	}
}

func (n Narwhal) Load(volume string, tarPath string) []string {
	id, err := uuid.NewUUID()
	if err != nil {
		return []string{err.Error()}
	}
	var errors []string
	name := "narwhal-mount" + "abc" + id.String()
	c := Container{name, n.quiet}

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

func (n Narwhal) Save(volume string, tarName string, path string) []string {
	id, err := uuid.NewUUID()
	if err != nil {
		return []string{err.Error()}
	}
	var errors []string

	zipped := tarName + ".tar.gz"
	name := "narwhal-mount" + "abc" + id.String()
	c := Container{name, n.quiet}

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

func (n Narwhal) KillAll() []string {

	containers, err := n.docker.ContainerIds()

	if err != "" {
		return []string{err}
	}

	if len(containers) == 0 {
		n.Print("No containers killed")
		return []string{}
	}

	n.Print("Killing all containers")

	containers = append([]string{"kill"}, containers...)

	kill := CreateCommand("docker", containers...)
	errs := kill.Run(n.quiet)
	if len(errs) != 0 {
		return errs
	}
	return []string{}

}

func (n Narwhal) RemoveAll() []string {
	containers, err := n.docker.AllContainerIds()

	if err != "" {
		return []string{err}
	}

	if len(containers) == 0 {
		n.Print("No containers removed")
		return []string{}
	}
	n.Print("Removing all containers")

	containers = append([]string{"rm"}, containers...)

	rm := CreateCommand("docker", containers...)
	errs := rm.Run(n.quiet)
	if len(errs) != 0 {
		return errs
	}
	return []string{}
}
