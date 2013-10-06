package isucontools

import (
	"strings"
)

func SqlEscape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}
