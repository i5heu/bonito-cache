package helper

import (
	"unicode"
)

func SanitizeMimeType(mime string) string {
	if mime == "" {
		return "application/octet-stream"
	}
	if len([]byte(mime)) > 256 {
		return "application/octet-stream"
	}

	// check if mime is only ascii
	for _, rune := range mime {
		if unicode.IsLetter(rune) == false && unicode.IsNumber(rune) == false && rune != '/' && rune != '.' && rune != '-' && rune != '_' {
			return "application/octet-stream"
		}
	}

	return mime
}
