package docker

import (
	"fmt"
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/kiringo/narwhal_lib/test_helper"
	"testing"
)

type DockerSuite struct {
	suite.Suite
	d *Docker
	f *test_helper.TestCommandFactory
}

func TestDocker(t *testing.T) {
	suite.Run(t, new(DockerSuite))
}

func (s *DockerSuite) SetupTest() {
	f := &test_helper.TestCommandFactory{[]string{}}
	d := &Docker{
		quiet: false,
		cmd:   f,
	}
	s.d = d
	s.f = f
}

func (s *DockerSuite) TearDownSuite() {
	clear := test_helper.HelpRun("docker", "rmi", "wew:tag")
	fmt.Println(clear)
}

func (s *DockerSuite) Test_Build() {
	assert := a.New(s.T())
	//test
	err := s.d.Build("../random", "do.ckerfile", "wew:tag")
	ret := test_helper.HelpRun("docker", "images", "--format", "{{.Repository}}:{{.Tag}}", "-f", "reference=wew")

	//assertions
	assert.Len(err, 0)
	assert.Equal(ret[0], "wew:tag")
	assert.Contains(s.f.Output, "Successfully tagged wew:tag")

}

func (s *DockerSuite) Test_Run() {
	assert := a.New(s.T())

	// test
	err := s.d.Run("wew:tag", "")
	assert.Empty(err)
	assert.Equal(s.f.Output[0], "BOOOOOO")
}
