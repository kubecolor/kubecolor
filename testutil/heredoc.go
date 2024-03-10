package testutil

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
)

var heredocReplacer = strings.NewReplacer(
	`\n`, "\n",
	`\r`, "\r",
	`\t`, "\t",
	`\e`, "\033", // ANSI escape character
)

func NewHereDoc(s string) string {
	return heredocReplacer.Replace(heredoc.Doc(s))
}

func NewHereDocf(s string, args ...interface{}) string {
	return heredocReplacer.Replace(heredoc.Docf(s, args...))
}
