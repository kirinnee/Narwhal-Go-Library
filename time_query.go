package narwhal_lib

import (
	"time"
)

type Image struct {
	id        string
	createdAt time.Time
}

type Images struct {
	image []Image
}

const (
	SINCE  = iota
	BEFORE = iota
	FROM   = iota
)

type Query struct {
}

//func (i *Images) TimeQuery(s []string) *Images {
//	for _,v := range s {
//
//	}
//	strings.Split(s, "=")
//
//
//}
