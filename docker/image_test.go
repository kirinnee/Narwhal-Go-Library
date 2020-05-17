package docker

import (
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/kiringo/narwhal_lib/test_helper"
	"sort"
	"testing"
)

type ImageSuite struct {
	suite.Suite
	d Docker
}

func TestImage(t *testing.T) {
	suite.Run(t, new(ImageSuite))

}

func (s *ImageSuite) SetupSuite() {
	test_helper.HelpRunQ("docker kill $(docker ps -q)")
	test_helper.HelpRunQ("docker rm $(docker ps -aq)")
	s.d = New(false, &test_helper.TestCommandFactory{Output: []string{}})
}

func (s *ImageSuite) TearDownSuite() {
	test_helper.HelpRunQ("docker kill $(docker ps -q)")
	test_helper.HelpRunQ("docker rm $(docker ps -aq)")
}

func (s *ImageSuite) TearDownTest() {
	test_helper.HelpRunQ(`docker rmi $(docker images --no-trunc -af "reference=narwhal/**/*" --format "{{.Repository}}:{{.Tag}}")`)
}

func (s *ImageSuite) Test_Image_labels() {

	//Setup
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first:1 --label a")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first:2 --label b")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first:3 --label c=1")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first:4 --label c=2")

	expect1 := []string{"narwhal/first:1"}
	expect2 := []string{"narwhal/first:2"}
	expect3 := []string{"narwhal/first:3", "narwhal/first:4"}
	expect4 := []string{"narwhal/first:3"}
	expect5 := make([]string, 0)

	assert := a.New(s.T())

	// subject
	image1, err1 := s.d.Images("", "a", []string{})
	image2, err2 := s.d.Images("", "b", []string{})
	image3, err3 := s.d.Images("", "c", []string{})
	image4, err4 := s.d.Images("", "c=1", []string{})
	image5, err5 := s.d.Images("", "c=3", []string{})
	//assert
	assert.Empty(err1)
	assert.Empty(err2)
	assert.Empty(err3)
	assert.Empty(err4)
	assert.Empty(err5)

	assert.ElementsMatch(expect1, toNames(image1))
	assert.ElementsMatch(expect2, toNames(image2))
	assert.ElementsMatch(expect3, toNames(image3))
	assert.ElementsMatch(expect4, toNames(image4))
	assert.ElementsMatch(expect5, toNames(image5))

}

func (s *ImageSuite) Test_Image_references() {

	//Setup
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first/first:1")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first/first:2")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first/second:1")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/first/second:2")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/second/first:3")
	test_helper.HelpRunQ("docker build ../labels --tag narwhal/second/first:4")

	expect1 := []string{"narwhal/first/first:1", "narwhal/first/first:2", "narwhal/first/second:1", "narwhal/first/second:2"}
	expect2 := []string{"narwhal/first/first:1", "narwhal/first/second:1"}
	expect3 := []string{"narwhal/second/first:3", "narwhal/second/first:4"}
	expect4 := []string{"narwhal/first/first:1", "narwhal/first/first:2", "narwhal/second/first:3", "narwhal/second/first:4"}
	expect5 := make([]string, 0)
	expect6 := []string{"narwhal/first/first:1", "narwhal/first/first:2", "narwhal/first/second:1", "narwhal/first/second:2"}

	assert := a.New(s.T())

	// subject
	image1, err1 := s.d.Images("", "", []string{"narwhal/first/*"})
	image2, err2 := s.d.Images("", "", []string{"narwhal/first/*:1"})
	image3, err3 := s.d.Images("", "", []string{"narwhal/second/*"})
	image4, err4 := s.d.Images("", "", []string{"narwhal/*/first:*"})
	image5, err5 := s.d.Images("", "", []string{"narwhal/narwhal/*"})
	image6, err6 := s.d.Images("", "", []string{"narwhal/first/first:*", "narwhal/first/second:*"})

	//assert
	assert.Empty(err1)
	assert.Empty(err2)
	assert.Empty(err3)
	assert.Empty(err4)
	assert.Empty(err5)
	assert.Empty(err6)

	assert.Equal(order(expect1), toNames(image1))
	assert.Equal(order(expect2), toNames(image2))
	assert.Equal(order(expect3), toNames(image3))
	assert.Equal(order(expect4), toNames(image4))
	assert.Equal(order(expect5), toNames(image5))
	assert.Equal(order(expect6), toNames(image6))

}

func order(in []string) []string {
	sort.Strings(in)
	return in
}

func toNames(images []Image) []string {
	in := make([]string, 0, len(images))
	for _, v := range images {
		in = append(in, v.Name)
	}
	sort.Strings(in)
	return in
}
