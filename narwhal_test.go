package narwhal_lib

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func helpRun(command string, arg ...string) []string {
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
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := helpRun("docker", "ps", "-q")
	if len(started) < 3 {
		t.Fail()
	}

	n.KillAll()
	left := helpRun("docker", "ps", "-q")
	fmt.Println("Left:", left)
	if len(left) != 0 {
		t.Fail()
	}
}

func TestNarwhal_StopAll(t *testing.T) {
	n := New(false)
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := helpRun("docker", "ps", "-q")
	if len(started) < 3 {
		t.Fail()
	}

	n.StopAll()

	left := helpRun("docker", "ps", "-q")
	fmt.Println("Left:", left)
	if len(left) != 0 {
		t.Fail()
	}
}

func TestNarwhal_RemoveAll(t *testing.T) {
	//setup
	n := New(false)
	helpRun("docker", "run", "hello-world")
	helpRun("docker", "run", "hello-world")
	helpRun("docker", "run", "hello-world")

	time.Sleep(1)
	started := helpRun("docker", "ps", "-aq")
	if len(started) < 3 {
		fmt.Print("not enough containers")
		t.Fail()
	}

	// test
	n.RemoveAll()

	left := helpRun("docker", "ps", "-aq")
	fmt.Println("Left:", left)
	if len(left) != 0 {
		fmt.Print("containers not removed")
		t.Fail()
	}

}

func TestNarwhal_DeployAuto(t *testing.T) {
	n := New(false)
	n.KillAll()

	n.DeployAuto("test-stack", "stack.yml", false)

	stack := helpRun("docker", "stack", "ls")
	if len(stack) != 2 {
		fmt.Println(stack, len(stack))
		t.Error("Incorrect number of stacks")
	}
	time.Sleep(time.Second * 10)
	stacks := helpRun("docker", "ps", "--format", "\"{{.Names}}\"")
	for _, v := range stacks {
		if !strings.HasPrefix(v, "\"test-stack_rocket.") {
			t.Error("Incorrect name: ", v)
		}
	}

	helpRun("docker", "swarm", "leave", "--force")
}

func TestNarwhal_Run(t *testing.T) {
	n := New(false)
	n.Run("random", "do.ckerfile", "sample:sample")
}
