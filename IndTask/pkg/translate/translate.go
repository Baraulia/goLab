package translate

import (
	"regexp"
)

func Translate(s string) string {
	cyrillic := s
	for k, v := range enTranslations {
		r, _ := regexp.Compile(v)
		cyrillic = r.ReplaceAllString(cyrillic, k)
	}
	return cyrillic
}
