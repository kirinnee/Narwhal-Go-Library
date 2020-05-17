package images

import (
	"gitlab.com/kiringo/narwhal_lib/docker"
	"time"
)

type Images []docker.Image
type Queries []TimeQuery

const (
	AFTER  = iota
	BEFORE = iota
)

type TimeQuery struct {
	query int
	time  time.Time
}

func (c Images) ProcessQuery(queries Queries) Images {
	ret := make([]docker.Image, 0, len(c))

OUTER:
	for _, v := range c {
		for _, q := range queries {
			if q.query == BEFORE && !v.CreatedAt.Before(q.time) {
				continue OUTER
			} else if q.query == AFTER && !v.CreatedAt.After(q.time) {
				continue OUTER
			}
		}
		ret = append(ret, v)
	}

	return ret
}
