package narwhal_lib

import (
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type NarwhalSuite struct {
	suite.Suite
	n       *Narwhal
	factory *TestCommandFactory
}

func TestNarwhal(t *testing.T) {
	suite.Run(t, new(NarwhalSuite))
}

func (s *NarwhalSuite) SetupTest() {
	factory := TestCommandFactory{}
	s.factory = &factory
	s.n = NewCustom(false, &factory)

}

func (s *NarwhalSuite) TearDownTest() {
	helpRunPrint(`docker swarm leave --force`)
	helpRunPrint(`docker kill $(docker ps -q)`)
	helpRunPrint(`docker rm $(docker ps -aq)`)
}

func (s *NarwhalSuite) Test_Save() {
	s.n.Save("cyanprint", "data", "./")
}

func (s *NarwhalSuite) Test_Load() {
	s.n.Load("ezvol", "./data.tar.gz")
}

func (s *NarwhalSuite) Test_KillAll() {
	// Setup
	assert := a.New(s.T())
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := helpRun("docker", "ps", "-q")
	assert.Len(started, 3)

	// Test
	s.n.KillAll()

	// Assert
	left := helpRun("docker", "ps", "-q")
	assert.Empty(left)

}

func (s *NarwhalSuite) Test_StopAll() {
	assert := a.New(s.T())

	// Setup
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	helpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := helpRun("docker", "ps", "-q")
	assert.Len(started, 3)

	// Test
	s.n.StopAll()

	// Assert
	left := helpRun("docker", "ps", "-q")
	assert.Empty(left)
}

func (s *NarwhalSuite) Test_RemoveAll() {
	//Setup
	assert := a.New(s.T())
	n := New(false)
	helpRun("docker", "run", "hello-world")
	helpRun("docker", "run", "hello-world")
	helpRun("docker", "run", "hello-world")

	time.Sleep(1)
	started := helpRun("docker", "ps", "-aq")
	assert.Len(started, 3)

	// test
	n.RemoveAll()

	left := helpRun("docker", "ps", "-aq")
	assert.Empty(left)
}

func (s *NarwhalSuite) Test_DeployAuto() {
	assert := a.New(s.T())

	// Test
	s.n.DeployAuto("test-stack", "stack.yml", false)
	stack := helpRun("docker", "stack", "ls")
	container := helpRun("docker", "ps", "--format", "\"{{.Names}}\"")
	time.Sleep(time.Second * 10)
	for i, v := range container {
		container[i] = string([]rune(v)[:18])
	}

	// Assert
	assert.Len(stack, 2)
	for _, v := range container {
		assert.Equal(v, `"test-stack_rocket.`)
	}

}

func (s *NarwhalSuite) Test_Run() {
	assert := a.New(s.T())
	s.n.Run("random", "do.ckerfile", "sample:sample")
	out := s.factory.output
	assert.Equal(out[len(out)-1], "BOOOOOO")

}
