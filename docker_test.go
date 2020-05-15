package narwhal_lib

import (
	"fmt"
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DockerSuite struct {
	suite.Suite
	d *Docker
	f *TestCommandFactory
}

func TestDocker(t *testing.T) {
	suite.Run(t, new(DockerSuite))
}

func (s *DockerSuite) SetupTest() {
	f := &TestCommandFactory{output: []string{}}
	d := &Docker{
		quiet: false,
		cmd:   f,
	}
	s.d = d
	s.f = f
}

func (s *DockerSuite) TearDownSuite() {
	clear := helpRun("docker", "rmi", "wew:tag")
	fmt.Println(clear)
}

func (s *DockerSuite) Test_Example() {
	assert := a.New(s.T())
	assert.Equal(5, 10-5)
}

func (s *DockerSuite) Test_Build() {
	assert := a.New(s.T())
	//test
	err := s.d.Build("random", "do.ckerfile", "wew:tag")
	ret := helpRun("docker", "images", "--format", "{{.Repository}}:{{.Tag}}", "-f", "reference=wew")

	//assertions
	assert.Len(err, 0)
	assert.Equal(ret[0], "wew:tag")
	assert.Contains(s.f.output, "Successfully tagged wew:tag")

}

func (s *DockerSuite) Test_Run() {
	assert := a.New(s.T())

	// test
	err := s.d.Run("wew:tag")
	assert.Empty(err)
	assert.Equal(s.f.output[0], "BOOOOOO")
}
