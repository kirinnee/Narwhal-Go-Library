package images

import (
	"errors"
	"github.com/araddon/dateparse"
	"gitlab.com/kiringo/narwhal_lib/docker"
	"strings"
	"time"
)

type RawQuery struct {
	query int
	time  string
}

type RawQueryParser struct {
	raw  []RawQuery
	done []TimeQuery
}

func (r *RawQueryParser) resolve() ([]TimeQuery, error) {
	if len(r.raw) == 0 {
		return r.done, nil
	} else {
		errs := make([]string, 0, len(r.raw))
		for _, v := range r.raw {
			errs = append(errs, v.time)
		}
		err := "Unknown queries: " + strings.Join(errs, ",")
		return nil, errors.New(err)
	}
}

func (r *RawQueryParser) parseTime() *RawQueryParser {
	raw, done := make([]RawQuery, 0), r.done

	for _, v := range r.raw {
		t, err := dateparse.ParseLocal(v.time)

		if err == nil {
			done = append(done, TimeQuery{
				query: v.query,
				time:  t,
			})
		} else {
			raw = append(raw, v)
		}
	}
	r.raw = raw
	r.done = done
	return r
}

func (r *RawQueryParser) parseDuration(t time.Time) *RawQueryParser {
	raw, done := make([]RawQuery, 0), r.done

	for _, v := range r.raw {
		duration, err := time.ParseDuration(v.time)
		if err == nil {
			done = append(done, TimeQuery{
				query: v.query,
				time:  t.Add(-duration),
			})
		} else {
			raw = append(raw, v)
		}
	}
	r.raw = raw
	r.done = done
	return r
}

func (r *RawQueryParser) parseImage(images []docker.Image) *RawQueryParser {
	raw, done := make([]RawQuery, 0), r.done
OUTER:
	for _, v := range r.raw {
		for _, i := range images {
			if i.Name == v.time {
				done = append(done, TimeQuery{
					query: v.query,
					time:  i.CreatedAt,
				})
				continue OUTER
			}
		}
		raw = append(raw, v)
	}
	r.raw = raw
	r.done = done
	return r
}

func IsTimeQuery(s string) bool {
	return strings.HasPrefix(s, "from=") ||
		strings.HasPrefix(s, "after=") ||
		strings.HasPrefix(s, "before=")
}

func New(s ...string) (*RawQueryParser, error) {

	all := make([][]RawQuery, 0, 10)
	for _, v := range s {
		raw, err := parseRawQuery(v)
		if err != nil {
			return nil, err
		}
		all = append(all, raw)
	}
	return newRawQueryParser(all...), nil

}

func newRawQueryParser(rq ...[]RawQuery) *RawQueryParser {
	raw := make([]RawQuery, 0, 100)
	for _, v := range rq {
		raw = append(raw, v...)
	}
	return &RawQueryParser{
		raw:  raw,
		done: []TimeQuery{},
	}
}

func (r *RawQueryParser) Parse(images []docker.Image) ([]TimeQuery, error) {
	r.parseImage(images)
	r.parseTime()
	r.parseDuration(time.Now())
	return r.resolve()
}

func parseRawQuery(s string) ([]RawQuery, error) {
	if !IsTimeQuery(s) {
		return nil, errors.New("not a time query")
	}
	if strings.HasPrefix(s, "from=") {
		q := strings.Split(s, " to=")
		from := strings.Replace(q[0], "from=", "", 1)
		to := q[1]
		return []RawQuery{{
			time:  from,
			query: AFTER,
		}, {
			time:  to,
			query: BEFORE,
		},
		}, nil
	} else if strings.HasPrefix(s, "after=") {
		return []RawQuery{{
			time:  strings.Split(s, "=")[1],
			query: AFTER,
		}}, nil
	} else {
		return []RawQuery{{
			time:  strings.Split(s, "=")[1],
			query: BEFORE,
		}}, nil
	}
}
