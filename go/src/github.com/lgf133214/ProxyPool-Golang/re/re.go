package re

import "regexp"

var(
	Compile1 = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+\\.\\d+:\\d+)")
	Compile2 = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+\\.\\d+)[^\\d]+(\\d+)")
)