package svgparser

import (
	"bytes"
)

func ParseTag(content string, start int) (string, bool, int) {
	buff := bytes.NewBufferString("")
	ending := false
	length := 0

	for i := start; i < len(content); i++ {
		char := content[i]

		if char == '<' {
			continue
		}

		if char == '/' && content[i-1] == '<' {
			ending = true
			continue
		}

		// Checks if the char is not a hyphen, lowercase or uppercase letter
		if char != '-' && (char < 97 || 122 < char) && (char < 65 || 90 < char) {
			length = i - start + 1
			break
		}

		buff.WriteByte(char)
	}

	for i := length + start - 1; i < len(content); i++ {
		if content[i] == '>' {
			length += (i - start - length) + 2
			break
		}
	}

	return buff.String(), ending, length
}
