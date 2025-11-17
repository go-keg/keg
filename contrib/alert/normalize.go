package alert

import "regexp"

// 归一化规则示例（可自行扩展）
var (
	ipRegex      = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	numRegex     = regexp.MustCompile(`\b\d+\b`)
	uuidRegex    = regexp.MustCompile(`[0-9a-fA-F-]{36}`)
	timestampReg = regexp.MustCompile(`\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}:\d{2}`)
)

func NormalizeError(msg string) string {
	msg = ipRegex.ReplaceAllString(msg, "<ip>")
	msg = uuidRegex.ReplaceAllString(msg, "<uuid>")
	msg = timestampReg.ReplaceAllString(msg, "<time>")
	msg = numRegex.ReplaceAllString(msg, "<num>")
	return msg
}