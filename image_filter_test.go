package narwhal_lib

import (
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sort"
	"testing"
)

type ImageFilterSuite struct {
	suite.Suite
}

func TestImageFilter(t *testing.T) {
	suite.Run(t, new(ImageFilterSuite))
}

func (s *ImageFilterSuite) Test_Catching_Single_Dangle() {
	assert := a.New(s.T())

	// test
	_, dangle1, _, _, _, _ := parseFilters("dangling=true", "label=b")
	_, dangle2, _, _, _, _ := parseFilters("dangling=false", "label=b")
	_, dangle3, _, _, _, _ := parseFilters("label=b")
	assert.Equal(dangle1, "true")
	assert.Equal(dangle2, "false")
	assert.Equal(dangle3, "")
}

func (s *ImageFilterSuite) Test_Catching_Multiple_Dangle() {
	assert := a.New(s.T())

	// test
	_, _, _, _, _, err := parseFilters("dangling=true", "dangling=false")
	assert.Contains(err, "multiple dangle filter")
}

func (s *ImageFilterSuite) Test_Catching_Single_Label() {
	assert := a.New(s.T())

	// test
	label1, _, _, _, _, _ := parseFilters("dangling=true", "label=b")
	label2, _, _, _, _, _ := parseFilters("dangling=true", "label=a")
	label3, _, _, _, _, _ := parseFilters("label=b=asd asd")
	label4, _, _, _, _, _ := parseFilters("ref=x")

	assert.Equal(label1, "b")
	assert.Equal(label2, "a")
	assert.Equal(label3, "b=asd asd")
	assert.Equal(label4, "")
}

func (s *ImageFilterSuite) Test_Catching_Multiple_Label() {
	assert := a.New(s.T())

	// test
	_, _, _, _, _, err := parseFilters("label=a", "label=b")
	assert.Contains(err, "multiple label filter")
}

func (s *ImageFilterSuite) Test_Catching_References() {
	assert := a.New(s.T())
	//expected
	e1 := []string{"a", "narwhal/*"}
	e2 := []string{"b", "test/**/*"}
	e3 := []string{"a"}
	e4 := []string{"a"}
	e5 := make([]string, 0)

	// test
	_, _, ref1, _, _, _ := parseFilters("ref=a", "ref=narwhal/*")
	_, _, ref2, _, _, _ := parseFilters("ref=b", "reference=test/**/*")
	_, _, ref3, _, _, _ := parseFilters("ref=a")
	_, _, ref4, _, _, _ := parseFilters("reference=a")
	_, _, ref5, _, _, _ := parseFilters("from=a")

	//Asserts
	assert.Equal(e1, ref1)
	assert.Equal(e2, ref2)
	assert.Equal(e3, ref3)
	assert.Equal(e4, ref4)
	assert.Equal(e5, ref5)

}

func (s *ImageFilterSuite) Test_Catching_TimeQuery() {
	assert := a.New(s.T())

	e1 := order([]string{"from=golang:latest to=golang:5"})
	e2 := order([]string{"before=6h2m"})
	e3 := order([]string{"after=2020/07/09 18:00"})
	e4 := order([]string{"after=1h2s", "before=rocket:rs"})
	e5 := order([]string{"from=2m to=2020-05-19", "after=alpine:3.9"})
	e6 := order([]string{})

	// test
	_, _, _, tq1, _, _ := parseFilters("from=golang:latest to=golang:5", "ref=a", "ref=narwhal/*")
	_, _, _, tq2, _, _ := parseFilters("before=6h2m", "ref=b", "reference=test/**/*")
	_, _, _, tq3, _, _ := parseFilters("after=2020/07/09 18:00", "ref=a")
	_, _, _, tq4, _, _ := parseFilters("after=1h2s", "before=rocket:rs", "reference=a")
	_, _, _, tq5, _, _ := parseFilters("from=2m to=2020-05-19", "after=alpine:3.9", "label=a")
	_, _, _, tq6, _, _ := parseFilters("label=a")

	assert.Equal(e1, order(tq1))
	assert.Equal(e2, order(tq2))
	assert.Equal(e3, order(tq3))
	assert.Equal(e4, order(tq4))
	assert.Equal(e5, order(tq5))
	assert.Equal(e6, order(tq6))

}

func (s *ImageFilterSuite) Test_Catching_Remainders() {
	assert := a.New(s.T())

	e1 := order([]string{"--format", "{{.Tag}}"})
	e2 := order([]string{"-a"})
	e3 := order([]string{"--all", "--digests"})
	e4 := order([]string{"--no-trunc"})
	e5 := order([]string{})
	e6 := order([]string{"-aq"})

	// test
	_, _, _, _, remain1, _ := parseFilters("--format", "{{.Tag}}", "from=golang:latest to=golang:5", "ref=a", "ref=narwhal/*")
	_, _, _, _, remain2, _ := parseFilters("before=6h2m", "ref=b", "-a", "reference=test/**/*")
	_, _, _, _, remain3, _ := parseFilters("--all", "--digests", "after=2020/07/09 18:00", "ref=a")
	_, _, _, _, remain4, _ := parseFilters("after=1h2s", "before=rocket:rs", "--no-trunc", "reference=a")
	_, _, _, _, remain5, _ := parseFilters("from=2m to=2020-05-19", "after=alpine:3.9", "label=a")
	_, _, _, _, remain6, _ := parseFilters("label=a", "-aq")

	assert.Equal(e1, order(remain1))
	assert.Equal(e2, order(remain2))
	assert.Equal(e3, order(remain3))
	assert.Equal(e4, order(remain4))
	assert.Equal(e5, order(remain5))
	assert.Equal(e6, order(remain6))

}

func order(in []string) []string {
	sort.Strings(in)
	return in
}
