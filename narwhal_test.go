package narwhal_lib

import (
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/kiringo/narwhal_lib/test_helper"
	"testing"
	"time"
)

type NarwhalSuite struct {
	suite.Suite
	n       *Narwhal
	factory *test_helper.TestCommandFactory
}

func TestNarwhal(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(NarwhalSuite))
}

func (s *NarwhalSuite) SetupTest() {
	factory := test_helper.TestCommandFactory{}
	s.factory = &factory
	s.n = NewCustom(false, &factory)

}

func (s *NarwhalSuite) TearDownTest() {
	test_helper.HelpRunPrint(`docker swarm leave --force`)
	test_helper.HelpRunPrint(`docker kill $(docker ps -q)`)
	test_helper.HelpRunPrint(`docker rm $(docker ps -aq)`)
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
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := test_helper.HelpRun("docker", "ps", "-q")
	assert.Len(started, 3)

	// Test
	s.n.KillAll()

	// Assert
	left := test_helper.HelpRun("docker", "ps", "-q")
	assert.Empty(left)

}

func (s *NarwhalSuite) Test_StopAll() {
	assert := a.New(s.T())

	// Setup
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	test_helper.HelpRun("docker", "run", "--rm", "-itd", "kirinnee/rocketrs:latest")
	started := test_helper.HelpRun("docker", "ps", "-q")
	assert.Len(started, 3)

	// Test
	s.n.StopAll()

	// Assert
	left := test_helper.HelpRun("docker", "ps", "-q")
	assert.Empty(left)
}

func (s *NarwhalSuite) Test_RemoveAll() {
	//Setup
	assert := a.New(s.T())
	n := New(false)
	test_helper.HelpRun("docker", "run", "hello-world")
	test_helper.HelpRun("docker", "run", "hello-world")
	test_helper.HelpRun("docker", "run", "hello-world")

	time.Sleep(1)
	started := test_helper.HelpRun("docker", "ps", "-aq")
	assert.Len(started, 3)

	// test
	n.RemoveAll()

	left := test_helper.HelpRun("docker", "ps", "-aq")
	assert.Empty(left)
}

func (s *NarwhalSuite) Test_DeployAuto() {
	assert := a.New(s.T())

	// Test
	s.n.DeployAuto("test-stack", "stack.yml", false)
	stack := test_helper.HelpRun("docker", "stack", "ls")
	container := test_helper.HelpRun("docker", "ps", "--format", "\"{{.Names}}\"")
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
	out := s.factory.Output
	assert.Equal(out[len(out)-1], "BOOOOOO")

}
