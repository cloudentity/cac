package diff

import (
	"fmt"
	"strings"
)

// https://github.com/google/go-cmp/issues/230#issuecomment-665750648
func colorize(in string) string {
	escapeCode := func(code int) string {
		return fmt.Sprintf("\x1b[%dm", code)
	}

	if in == "" {
		return ""
	}

	ss := strings.Split(in, "\n")
	for i, s := range ss {
		switch {
		case strings.HasPrefix(s, "-"):
			ss[i] = escapeCode(31) + s + escapeCode(0)
		case strings.HasPrefix(s, "+"):
			ss[i] = escapeCode(32) + s + escapeCode(0)
		}
	}

	return strings.Join(ss, "\n")
}
