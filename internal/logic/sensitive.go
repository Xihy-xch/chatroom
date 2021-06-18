package logic

import (
	"strings"

	"github.com/Xihy-xch/tcp-chatroom/global"
)

func FilterSensitive(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}
