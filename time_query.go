package narwhal_lib

import (
	"time"
)

type Image struct {
	name      string
	id        string
	createdAt time.Time
}

type Images struct {
	image []Image
}

const (
	SINCE  = iota
	BEFORE = iota
)

type Query struct {
	query int
	time  time.Time
}

type RawQuery struct {
	query int
	time  string
}

//func parseRawQuery(s string) ([]RawQuery, error) {
//
//}
//
//func (i *Images) TimeQuery(s string) *Images {
//	strings.Split(s, "=")
//}
