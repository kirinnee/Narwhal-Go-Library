package narwhal_lib

import (
	"fmt"
	"testing"
)

func HelpRun(command string, arg ...string) []string {
	ret := make([]string, 0, 10)
	CreateCommand(command, arg...).CustomRun(false, func(s string) {
		ret = append(ret, s)
	}, func(s string) {
		ret = append(ret, s)
	})
	return ret
}

func TestNarwhal_Save(t *testing.T) {
	n := New(false)
	n.Save("cyanprint", "data", "./")
}

func TestNarwhal_Load(t *testing.T) {
	n := New(false)
	n.Load("ezvol", "./data.tar.gz")
}

func TestNarwhal_KillAll(t *testing.T) {
	n := New(false)
	HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := HelpRun("docker", "ps", "-q")
	if len(started) < 3 {
		t.Fail()
	}

	n.KillAll()
	left := HelpRun("docker", "ps", "-q")
	fmt.Println("Left:", left)
	if len(left) != 0 {
		t.Fail()
	}
}
