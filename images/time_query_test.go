package images

import (
	"github.com/araddon/dateparse"
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TimeQuerySuite struct {
	suite.Suite
}

func TestTimeQuery(t *testing.T) {
	suite.Run(t, new(TimeQuerySuite))
}

func (s *TimeQuerySuite) Test_ProcessQueries() {
	assert := a.New(s.T())

	//setup
	date1, err1 := dateparse.ParseLocal("2020-05-20 17:50:15") //after
	date2, err2 := dateparse.ParseLocal("2020-05-20 18:12:30") //after
	date3, err3 := dateparse.ParseLocal("2020-05-20 19:30:45") //before

	// image dates
	iDate1, iErr1 := dateparse.ParseLocal("2020-05-21 17:50:15")   // too late
	iDate2, iErr2 := dateparse.ParseLocal("2020-04-21 17:50:02")   // too early
	iDate3, iErr3 := dateparse.ParseLocal("2020-05-20 17:50:56")   // too early
	iDate4, iErr4 := dateparse.ParseLocal("2020-05-20 18:20:25")   // stay
	iDate5, iErr5 := dateparse.ParseLocal("2020-05-20 18:57:16")   // stay
	iDate6, iErr6 := dateparse.ParseLocal("2020-05-30 19:50:50")   // too late
	iDate7, iErr7 := dateparse.ParseLocal("2021-06-20 18:20:50")   // too late
	iDate8, iErr8 := dateparse.ParseLocal("2020-05-20 19:00:00")   // stay
	iDate9, iErr9 := dateparse.ParseLocal("2020-05-20 18:12:33")   // stay
	iDate10, iErr10 := dateparse.ParseLocal("2020-05-20 08:24:20") // too early
	iDate11, iErr11 := dateparse.ParseLocal("2013-08-06 13:16:20") // too early
	iDate12, iErr12 := dateparse.ParseLocal("2020-08-06 23:23:24") // too late

	eDate1, eErr1 := dateparse.ParseLocal("2020-05-20 18:20:25")
	eDate2, eErr2 := dateparse.ParseLocal("2020-05-20 18:57:16")
	eDate3, eErr3 := dateparse.ParseLocal("2020-05-20 19:00:00")
	eDate4, eErr4 := dateparse.ParseLocal("2020-05-20 18:12:33")

	// expected
	expected := Images{

		{
			Name:      "rocker:rs",
			Id:        "3",
			CreatedAt: eDate1,
		},
		{
			Name:      "goatling:go",
			Id:        "4",
			CreatedAt: eDate2,
		},

		{
			Name:      "node:express",
			Id:        "7",
			CreatedAt: eDate3,
		},
		{
			Name:      "postgres:12",
			Id:        "8",
			CreatedAt: eDate4,
		},
	}

	// Subjects
	queries := Queries{
		{
			query: BEFORE,
			time:  date3,
		},
		{
			query: AFTER,
			time:  date2,
		}, {
			query: AFTER,
			time:  date1,
		},
	}
	images := Images{
		{
			Name:      "image:1",
			Id:        "0",
			CreatedAt: iDate1,
		},
		{
			Name:      "image:2",
			Id:        "1",
			CreatedAt: iDate2,
		},
		{
			Name:      "image:3",
			Id:        "2",
			CreatedAt: iDate3,
		},
		{
			Name:      "rocker:rs",
			Id:        "3",
			CreatedAt: iDate4,
		},
		{
			Name:      "goatling:go",
			Id:        "4",
			CreatedAt: iDate5,
		},
		{
			Name:      "ruby:ror",
			Id:        "5",
			CreatedAt: iDate6,
		},
		{
			Name:      "rust:actix",
			Id:        "6",
			CreatedAt: iDate7,
		},
		{
			Name:      "node:express",
			Id:        "7",
			CreatedAt: iDate8,
		},
		{
			Name:      "postgres:12",
			Id:        "8",
			CreatedAt: iDate9,
		},

		{
			Name:      "alpine:13.9",
			Id:        "9",
			CreatedAt: iDate10,
		},

		{
			Name:      "ubuntu:latest",
			Id:        "10",
			CreatedAt: iDate11,
		},
		{
			Name:      "alpine:5",
			Id:        "11",
			CreatedAt: iDate12,
		},
	}

	actual := images.ProcessQuery(queries)

	// Asserts
	assertNil(assert, err1, err2, err3, eErr1, eErr2, eErr3, eErr4, iErr1, iErr2, iErr3, iErr4, iErr5, iErr6, iErr7, iErr8, iErr9, iErr10, iErr11, iErr12)
	assert.Equal(expected, actual)

}

func assertNil(assert *a.Assertions, err ...error) {
	for _, e := range err {
		assert.Nil(e)
	}
}
