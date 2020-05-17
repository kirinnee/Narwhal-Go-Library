package narwhal_lib

import (
	"gitlab.com/kiringo/narwhal_lib/images"
	"strings"
)

const DANGLE = "dangling="
const REF1 = "ref="
const REF2 = "reference="
const LABEL = "label="

func trim(word, trim string) string {
	return string([]rune(word)[len(trim):])
}

func parseFilters(filters ...string) (label, dangle string, ref, tq, remaining, err []string) {
	label, dangle = "", ""
	ref, tq, remaining, err = make([]string, 0, 20), make([]string, 0, 20), make([]string, 0, 20), make([]string, 0, 20)

	for _, v := range filters {
		if strings.HasPrefix(v, DANGLE) {
			if dangle == "" {
				dangle = trim(v, DANGLE)
			} else {
				err = append(err, "multiple dangle filter")
			}
		} else if strings.HasPrefix(v, LABEL) {
			if label == "" {
				label = trim(v, LABEL)
			} else {
				err = append(err, "multiple label filter")
			}
		} else if strings.HasPrefix(v, REF1) {
			ref = append(ref, trim(v, REF1))
		} else if strings.HasPrefix(v, REF2) {
			ref = append(ref, trim(v, REF2))
		} else if images.IsTimeQuery(v) {
			tq = append(tq, v)
		} else {
			remaining = append(remaining, v)
		}

	}
	return

}
