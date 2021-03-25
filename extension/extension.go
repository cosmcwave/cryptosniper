package extension

import "strings"

func Parse(ext string) (string, string) {
	extParsed := strings.Split(ext, ":")
	return extParsed[0], extParsed[1]
}