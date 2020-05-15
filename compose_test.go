package narwhal_lib

import (
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func breakDown() (map[string][]byte, error) {
	b, err := ioutil.ReadFile("./compose_test/sample.yml")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string][]byte)
	for k, v := range data {

		out, err := yaml.Marshal(&v)
		if err != nil {
			return nil, err
		}
		ret[k] = out
	}
	return ret, nil
}

type ComposeSuite struct {
	suite.Suite
	m *map[string][]byte
}

func TestCompose(t *testing.T) {
	suite.Run(t, new(ComposeSuite))
}

func (s *ComposeSuite) SetupSuite() {
	m, err := breakDown()
	if err != nil {
		panic("failed ")
	}
	s.m = &m
}

//func (s *ComposeSuite) TearDownSuite() {
//
//}

func (s *ComposeSuite) Test_Example() {
	assert := a.New(s.T())
	assert.Equal(5, 10-5)

}

func (s *ComposeSuite) Test_ParseConfig_both_context_and_file() {
	assert := a.New(s.T())
	m := *s.m

	// test
	subj := m["rocket"]
	expected := Builds{
		Context: "random",
		File:    "df",
	}
	actual, err := parseConfig(subj)

	// assert
	assert.Nil(err)
	assert.Equal(actual, expected)

}

func (s *ComposeSuite) Test_ParseConfig_string() {
	assert := a.New(s.T())
	m := *s.m

	// test
	subj := m["golang"]
	expected := Builds{
		Context: "welp",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)

	// assert
	assert.Nil(err)
	assert.Equal(actual, expected)

}

func (s *ComposeSuite) TestParseConfig_only_file() {
	assert := a.New(s.T())
	m := *s.m

	// test
	subj := m["dotnet"]

	expected := Builds{
		Context: ".",
		File:    "RandomFile",
	}
	actual, err := parseConfig(subj)

	// assert
	assert.Nil(err)
	assert.Equal(actual, expected)
}

func (s *ComposeSuite) TestParseConfig_only_context() {
	assert := a.New(s.T())
	m := *s.m

	// test
	subj := m["node"]

	expected := Builds{
		Context: "golang",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)

	// assert
	assert.Nil(err)
	assert.Equal(actual, expected)

}

func (s *ComposeSuite) TestParseConfig_empty() {
	assert := a.New(s.T())
	m := *s.m

	// test
	subj := m["ror"]

	expected := Builds{
		Context: ".",
		File:    "Dockerfile",
	}
	actual, err := parseConfig(subj)

	// assert
	assert.Nil(err)
	assert.Equal(actual, expected)

}
