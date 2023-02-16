package utils

import "github.com/microcosm-cc/bluemonday"

var policy = bluemonday.UGCPolicy()

func SanitizeHtml(input string) string {
	return policy.Sanitize(input)
}
