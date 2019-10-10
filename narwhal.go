package narwhal_lib

import (
	"github.com/google/uuid"
	"log"
)

type Narwhal struct {
	quiet bool
}

func (n Narwhal) Save(volume string, tarName string, path string) []string {
	id, err := uuid.NewUUID()
	if err != nil {
		return []string{err.Error()}
	}
	var errors []string
	name := "narwhal-mount" + "abc" + id.String()
	zipped := tarName + ".tar.gz"
	c := Container{name, n.quiet}

	log.Println("Creating container to connect to volume...")
	s := c.Start("alpine", volume, "/home/data", "dt")
	if s != "" {
		return []string{s}
	}
	log.Println("Container Created")

	log.Println("Zipping volume...")
	s = c.Exec("/home", "tar","-czf",zipped, "data")
	if s != ""{
		errors = append(errors, s)
	}
	log.Println("Volume Zipped!")

	if len(errors) == 0{
		log.Println("Copying to host...")
		s = c.Copy(false, "/home/" + zipped, path)
		if s != ""{
			errors = append(errors, s)
		}
		log.Println("Done copying!")
	}

	log.Println("Killing Container...")
	s = c.Kill()
	if s != "" {
		errors = append(errors, s)
		log.Println("Container failed to be killed!")
	}else {
		log.Println("Container Killed...")
	}

	log.Println("Removing Container...")
	s = c.Remove()
	if s != "" {
		errors = append(errors, s)
		log.Println("Container failed to be removed!")
	}else {
		log.Println("Container remove...")
	}

	return errors

}
