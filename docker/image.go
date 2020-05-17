package docker

import (
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

type (
	Image struct {
		Name      string
		Id        string
		CreatedAt time.Time
	}
)

func (d Docker) Images(dangling string, label string, references []string) (image []Image, err []string) {
	defer func() {
		if r := recover(); r != nil {
			err = []string{r.(error).Error()}
		}
	}()

	args := []string{
		"images", "--format", `{{.ID}}@{{.Repository}}:{{.Tag}}@{{.CreatedAt}}`,
	}
	if dangling != "" {
		if dangling == "true" {
			args = append(args, []string{"-f", `dangling=true`}...)
		} else if dangling == "false" {
			args = append(args, []string{"-f", `dangling=false`}...)
		}
	}

	if label != "" {
		args = append(args, []string{"-f", `label=` + label}...)

	}

	for _, v := range references {
		args = append(args, []string{"-f", `reference=` + v}...)

	}

	cmd := d.cmd.Create("docker", args...)

	images := make([]string, 0, 10)
	image = make([]Image, 0, 10)
	err = make([]string, 0, 10)
	err2 := cmd.CustomRun(func(s string) {
		images = append(images, s)
	}, func(s string) {
		err = append(err, s)
	})
	if err2 != "" {
		err = append(err, err2)
	}
	if len(err) != 0 {
		return
	}
	for _, v := range images {
		frag := strings.Split(v, "@")
		date, e := dateparse.ParseLocal(frag[2])
		if e != nil {
			return nil, []string{e.Error()}
		}
		image = append(image, Image{
			Name:      frag[1],
			Id:        frag[0],
			CreatedAt: date,
		})
	}
	return image, nil
}
