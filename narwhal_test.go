package narwhal_lib

import (
	"fmt"
	"strings"
	"testing"
	"time"
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

func TestNarwhal_RemoveAll(t *testing.T) {
	//setup
	n := New(false)
	HelpRun("docker", "run", "hello-world")
	HelpRun("docker", "run", "hello-world")
	HelpRun("docker", "run", "hello-world")

	time.Sleep(1)
	started := HelpRun("docker", "ps", "-aq")
	if len(started) < 3 {
		fmt.Print("not enough containers")
		t.Fail()
	}

	// test
	n.RemoveAll()

	left := HelpRun("docker", "ps", "-aq")
	fmt.Println("Left:", left)
	if len(left) != 0 {
		fmt.Print("containers not removed")
		t.Fail()
	}

}

func TestNarwhal_DeployAuto(t *testing.T) {
	n := New(false)
	n.KillAll()

	n.DeployAuto("test-stack", "stack.yml")

	stack := HelpRun("docker", "stack", "ls")
	if len(stack) != 2 {
		fmt.Println(stack, len(stack))
		t.Error("Incorrect number of stacks")
	}
	time.Sleep(time.Second * 10)
	stacks := HelpRun("docker", "ps", "--format", "\"{{.Names}}\"")
	for _, v := range stacks {
		if !strings.HasPrefix(v, "\"test-stack_web_app.") {
			t.Error("Incorrect name: ", v)
		}
	}

}
