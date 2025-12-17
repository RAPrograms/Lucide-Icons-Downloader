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

func FormatSVG(content string) string {
	//Getting root node
	tag, _, node_length := ParseTag(content, 0)

	//Adding root
	hierarchy := append(make([]string, 0), tag)
	output := content[0 : node_length-1]

	for i := node_length; i < len(content); i++ {
		tag, node_end, node_length := ParseTag(content, i)

		// Getting parent node and current node
		parent_tag := hierarchy[len(hierarchy)-1]
		node := content[max(0, i-1) : i+node_length-1]

		// Ensuring next loop is the next tag (if exists)
		i += node_length - 1

		// Checks for nodes without children
		if node_end && parent_tag == tag {
			output += node
			continue
		}

		// Adding or removing node, with children, to hieracrchy
		if node_end {
			hierarchy = hierarchy[:len(hierarchy)-1]
		} else {
			hierarchy = append(hierarchy, tag)
		}

		//Adds new line
		output += "\n"

		//Adds indentation
		for i := 0; i < len(hierarchy)-1; i++ {
			output += "\t"
		}

		output += node
	}

	return output
}
