package util

import (
	"regexp"
	"strings"
)

func GenerateBaseSlug(text string) string {
	text = strings.ToLower(text)

	reg, _ := regexp.Compile(`[^a-z0-9\\s]+`)
	text = reg.ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, " ", "-")

	return text
}
