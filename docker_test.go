package narwhal_lib

import (
	"fmt"
	"testing"
)

func TestDocker_Build(t *testing.T) {
	d := Docker{quiet: false}
	err := d.Build("random", "do.ckerfile", "wew:tag")
	if len(err) > 0 {
		t.Error(err)
	}
	ret := helpRun("docker", "images", "--format", "{{.Repository}}:{{.Tag}}", "-f", "reference=wew")
	if ret[0] != "wew:tag" {
		t.Error("Image not built", ret)
	}
}

func TestDocker_Run(t *testing.T) {
	d := Docker{quiet: false}

	err := d.Run("wew:tag")
	if len(err) > 0 {
		t.Error(err)
	}
	clear := helpRun("docker", "rmi", "wew:tag")
	fmt.Println(clear)
}
