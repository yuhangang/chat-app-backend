package output_processing

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

// ExtractHTMLNodes extracts all HTML nodes from content
func ExtractHTMLNodes(content string) ([]string, error) {
	var htmlBlocks []string

	// Create a tokenizer
	tokenizer := html.NewTokenizer(strings.NewReader(content))

	var buffer bytes.Buffer
	depth := 0
	var currentTag string

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of the document or error
			return htmlBlocks, nil

		case html.StartTagToken:
			tag, hasAttr := tokenizer.TagName()
			tagName := string(tag)

			if depth == 0 {
				currentTag = tagName
				buffer.Reset()
			}

			buffer.WriteString("<" + tagName)

			// Add attributes
			if hasAttr {
				for {
					key, val, more := tokenizer.TagAttr()
					buffer.WriteString(" " + string(key) + "=\"" + string(val) + "\"")
					if !more {
						break
					}
				}
			}

			buffer.WriteString(">")
			depth++

		case html.EndTagToken:
			tag, _ := tokenizer.TagName()
			tagName := string(tag)

			buffer.WriteString("</" + tagName + ">")
			depth--

			if depth == 0 && tagName == currentTag {
				htmlBlocks = append(htmlBlocks, buffer.String())
			}

		case html.SelfClosingTagToken:
			tag, hasAttr := tokenizer.TagName()
			tagName := string(tag)

			var selfClosingTag strings.Builder
			selfClosingTag.WriteString("<" + tagName)

			// Add attributes
			if hasAttr {
				for {
					key, val, more := tokenizer.TagAttr()
					selfClosingTag.WriteString(" " + string(key) + "=\"" + string(val) + "\"")
					if !more {
						break
					}
				}
			}

			selfClosingTag.WriteString("/>")
			htmlBlocks = append(htmlBlocks, selfClosingTag.String())

		case html.TextToken:
			text := tokenizer.Text()
			buffer.Write(text)
		}
	}
}
