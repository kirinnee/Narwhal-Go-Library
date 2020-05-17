package images

import (
	"github.com/araddon/dateparse"
	a "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/kiringo/narwhal_lib/docker"
	"testing"
)

type RawQuerySuite struct {
	suite.Suite
}

func TestRawTimeQuery(t *testing.T) {
	suite.Run(t, new(RawQuerySuite))
}

func (s *RawQuerySuite) SetupSuite() {

}

func (s *RawQuerySuite) TearDownSuite() {

}

func (s *RawQuerySuite) Test_ParseRawQuery_After_Relative() {
	assert := a.New(s.T())

	raw, err := parseRawQuery("after=6m2s")
	expected := newRawQueryParser([]RawQuery{{query: AFTER, time: "6m2s"}})

	assert.Nil(err)
	assert.Equal(raw, expected)
}

func (s *RawQuerySuite) Test_ParseRawQuery_After_Absolute() {
	assert := a.New(s.T())

	raw1, err1 := parseRawQuery("after=2020/05/07")
	expected1 := newRawQueryParser([]RawQuery{{query: AFTER, time: "2020/05/07"}})

	raw2, err2 := parseRawQuery("after=2020-01-03 09:21:37 +0800 +08")
	expected2 := newRawQueryParser([]RawQuery{{query: AFTER, time: "2020-01-03 09:21:37 +0800 +08"}})

	assert.Nil(err1)
	assert.Nil(err2)
	assert.Equal(raw1, expected1)
	assert.Equal(raw2, expected2)
}

func (s *RawQuerySuite) Test_ParseRawQuery_After_Image() {
	assert := a.New(s.T())

	raw1, err1 := parseRawQuery("after=alpine:latest")
	expected1 := newRawQueryParser([]RawQuery{{query: AFTER, time: "alpine:latest"}})

	raw2, err2 := parseRawQuery("after=ubuntu:tag")
	expected2 := newRawQueryParser([]RawQuery{{query: AFTER, time: "ubuntu:tag"}})

	assert.Nil(err1)
	assert.Nil(err2)
	assert.Equal(raw1, expected1)
	assert.Equal(raw2, expected2)
}

func (s *RawQuerySuite) Test_ParseRawQuery_Before_Relative() {
	assert := a.New(s.T())

	raw, err := parseRawQuery("before=6m2s")
	expected := newRawQueryParser([]RawQuery{{query: BEFORE, time: "6m2s"}})

	assert.Nil(err)
	assert.Equal(raw, expected)
}

func (s *RawQuerySuite) Test_ParseRawQuery_Before_Absolute() {
	assert := a.New(s.T())

	raw1, err1 := parseRawQuery("before=2020/05/07")
	expected1 := newRawQueryParser([]RawQuery{{query: BEFORE, time: "2020/05/07"}})

	raw2, err2 := parseRawQuery("before=2020-01-03 09:21:37 +0800 +08")
	expected2 := newRawQueryParser([]RawQuery{{query: BEFORE, time: "2020-01-03 09:21:37 +0800 +08"}})

	assert.Nil(err1)
	assert.Nil(err2)
	assert.Equal(raw1, expected1)
	assert.Equal(raw2, expected2)
}

func (s *RawQuerySuite) Test_ParseRawQuery_Before_Image() {
	assert := a.New(s.T())

	raw1, err1 := parseRawQuery("before=alpine:latest")
	expected1 := newRawQueryParser([]RawQuery{{query: BEFORE, time: "alpine:latest"}})

	raw2, err2 := parseRawQuery("before=ubuntu:tag")
	expected2 := newRawQueryParser([]RawQuery{{query: BEFORE, time: "ubuntu:tag"}})

	assert.Nil(err1)
	assert.Nil(err2)
	assert.Equal(raw1, expected1)
	assert.Equal(raw2, expected2)
}

func (s *RawQuerySuite) Test_ParseRawQuery_From() {
	assert := a.New(s.T())

	raw, err := parseRawQuery("from=2020-01-03 09:21:37 +0800 +08 to=image:latest")

	expected := newRawQueryParser([]RawQuery{{
		query: AFTER,
		time:  "2020-01-03 09:21:37 +0800 +08",
	}, {
		query: BEFORE,
		time:  "image:latest",
	}})

	assert.Nil(err)
	assert.Equal(raw, expected)
}

func (s *RawQuerySuite) Test_ParseRaw_Query_Reject_Non_string_Query() {
	assert := a.New(s.T())

	_, err1 := parseRawQuery("f=image:latest")
	_, err2 := parseRawQuery("other=image:latest")
	_, err3 := parseRawQuery("since=image:latest")
	assert.Error(err1)
	assert.Error(err2)
	assert.Error(err3)

}

func (s *RawQuerySuite) Test_ParseImage() {
	assert := a.New(s.T())

	// setup subject
	date1, err1 := dateparse.ParseLocal("2020-01-02")
	date2, err2 := dateparse.ParseLocal("2020-01-05 18:00")
	date3, err3 := dateparse.ParseLocal("2020-12-06 06:15:22")
	date4, err4 := dateparse.ParseLocal("2019-09-12 06:15:22")
	images := []docker.Image{
		{
			Name:      "postgres:13.0",
			Id:        "id1",
			CreatedAt: date1,
		},
		{
			Name:      "rocket:latest",
			Id:        "id1",
			CreatedAt: date2,
		},
		{
			Name:      "golang:latest",
			Id:        "id2",
			CreatedAt: date3,
		},
	}

	subj := &RawQueryParser{
		raw: []RawQuery{
			{
				query: BEFORE,
				time:  "golang:latest",
			},
			{
				query: AFTER,
				time:  "2020-01-05",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
			{
				query: AFTER,
				time:  "rocket:latest",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  date4,
			},
		},
	}

	// set up expected
	eDate2, err5 := dateparse.ParseLocal("2020-01-05 18:00")
	eDate3, err6 := dateparse.ParseLocal("2020-12-06 06:15:22")
	eDate4, err7 := dateparse.ParseLocal("2019-09-12 06:15:22")
	expected := &RawQueryParser{
		raw: []RawQuery{
			{
				query: AFTER,
				time:  "2020-01-05",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  eDate4,
			},
			{
				query: BEFORE,
				time:  eDate3,
			},
			{
				query: AFTER,
				time:  eDate2,
			},
		},
	}

	// test
	actual := subj.parseImage(images)

	//assert
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Nil(err4)
	assert.Nil(err5)
	assert.Nil(err6)
	assert.Nil(err7)
	assert.Equal(expected, actual)

}

func (s *RawQuerySuite) Test_ParseDuration() {
	assert := a.New(s.T())
	// setup subject

	rawDate, err0 := dateparse.ParseAny("2010-07-26 10:15:40")
	subj := &RawQueryParser{
		raw: []RawQuery{
			{
				query: BEFORE,
				time:  "rocket:latest",
			},
			{
				query: BEFORE,
				time:  "2h6m",
			},
			{
				query: AFTER,
				time:  "2020-01-05",
			},
			{
				query: BEFORE,
				time:  "20s",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
			{
				query: AFTER,
				time:  "2h15s",
			},
			{
				query: AFTER,
				time:  "2m20s",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  rawDate,
			},
		},
	}

	// set up expected
	date1, err1 := dateparse.ParseLocal("2020-05-16 05:58:45")
	date2, err2 := dateparse.ParseLocal("2020-05-16 08:04:25")
	date3, err3 := dateparse.ParseLocal("2020-05-16 06:04:30")
	date4, err4 := dateparse.ParseLocal("2020-05-16 08:02:25")
	eRawDate, err5 := dateparse.ParseAny("2010-07-26 10:15:40")
	expected := &RawQueryParser{
		raw: []RawQuery{
			{
				query: BEFORE,
				time:  "rocket:latest",
			},
			{
				query: AFTER,
				time:  "2020-01-05",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  eRawDate,
			},
			{
				query: BEFORE,
				time:  date1,
			},
			{
				query: BEFORE,
				time:  date2,
			},
			{
				query: AFTER,
				time:  date3,
			},
			{
				query: AFTER,
				time:  date4,
			},
		},
	}

	// test
	testTime, err := dateparse.ParseLocal("2020-05-16 08:04:45")
	actual := subj.parseDuration(testTime)

	//assert
	assert.Nil(err)
	assert.Nil(err0)
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Nil(err4)
	assert.Nil(err5)

	assert.Equal(expected, actual)

}

func (s *RawQuerySuite) Test_ParseTime() {
	assert := a.New(s.T())
	// setup subject

	rawDate, err0 := dateparse.ParseAny("2010-07-26 10:15:40")
	subj := &RawQueryParser{
		raw: []RawQuery{
			{
				query: BEFORE,
				time:  "rocket:latest",
			},
			{
				query: BEFORE,
				time:  "2h6m",
			},
			{
				query: AFTER,
				time:  "2020-01-05",
			},
			{
				query: BEFORE,
				time:  "Mon Jan  2 15:04:05 2006",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
			{
				query: AFTER,
				time:  "03 February 2013",
			},
			{
				query: AFTER,
				time:  "8/8/1965 01:00 PM ",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  rawDate,
			},
		},
	}

	// set up expected
	date1, err1 := dateparse.ParseLocal("2020-01-05 00:00:00")
	date2, err2 := dateparse.ParseLocal("2006-01-02 15:04:05")
	date3, err3 := dateparse.ParseLocal("2013-02-03 00:00:00")
	date4, err4 := dateparse.ParseLocal("1965-08-08 13:00:00")
	eRawDate, err5 := dateparse.ParseAny("2010-07-26 10:15:40")
	expected := &RawQueryParser{
		raw: []RawQuery{
			{
				query: BEFORE,
				time:  "rocket:latest",
			},
			{
				query: BEFORE,
				time:  "2h6m",
			},
			{
				query: BEFORE,
				time:  "postgres:latest",
			},
		},
		done: []TimeQuery{
			{
				query: AFTER,
				time:  eRawDate,
			},
			{
				query: AFTER,
				time:  date1,
			},
			{
				query: BEFORE,
				time:  date2,
			},
			{
				query: AFTER,
				time:  date3,
			},
			{
				query: AFTER,
				time:  date4,
			},
		},
	}

	// test
	actual := subj.parseTime()

	//assert
	assert.Nil(err0)
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Nil(err4)
	assert.Nil(err5)

	assert.Equal(expected, actual)

}
