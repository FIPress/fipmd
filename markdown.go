package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil"
)

/*
var tags = map[string]struct{} {
	"a":			struct{}{},
	"abbr":			struct{}{},
	"address":		struct{}{},
	"area":			struct{}{},
	"article":    	struct{}{},
	"aside":      	struct{}{},
	"audio":		struct{}{},
	"b":			struct{}{},
	"base":			struct{}{},
	"bdi":			struct{}{},
	"bdo":			struct{}{},
	"blockquote": 	struct{}{},
	"body":			struct{}{},
	"canvas":     struct{}{},
	"center":		struct{}{},
	"del":        struct{}{},
	"div":        struct{}{},
	"dl":         struct{}{},
	"dt":         struct{}{},
	"dd":         struct{}{},
	"fieldset":   struct{}{},
	"figcaption": struct{}{},
	"figure":     struct{}{},

	"footer":     struct{}{},
	"form":       struct{}{},
	"h1":         struct{}{},
	"h2":         struct{}{},
	"h3":         struct{}{},
	"h4":         struct{}{},
	"h5":         struct{}{},
	"h6":         struct{}{},
	"head":			struct{}{},
	"header":     struct{}{},
	"hgroup":     struct{}{},
	"html":			struct{}{},
	"iframe":     struct{}{},
	"ins":        struct{}{},
	"li":			struct{}{},
	"link":			struct{}{},
	"math":       struct{}{},
	"meta":			struct{}{},
	"noscript":   struct{}{},
	"ol":         struct{}{},
	"output":     struct{}{},
	"p":          struct{}{},
	"pre":        struct{}{},
	"progress":   struct{}{},
	"section":    struct{}{},
	"script":		struct{}{},
	"style":			struct{}{},
	"table":      struct{}{},
	"title":			struct{}{},
	"ul":         struct{}{},
	"video":      struct{}{},
}*/

type markdown interface {
	//write html to buffer
	Html() string
	Text() string
	html(buf *bytes.Buffer)
	inText() []byte
	execHooks(parser *Parser)
}

func writeSpan(tag string, content []byte, buf *bytes.Buffer) {
	buf.Write(startTag(tag))
	HtmlEscapeToBuffer(content, buf)
	buf.Write(endTag(tag))
}

func writeSpanLiterally(tag string, content []byte, buf *bytes.Buffer) {
	buf.Write(startTag(tag))
	HtmlEscapeLiterallyToBuffer(content, buf)
	buf.Write(endTag(tag))
}

func startTag(tag string) []byte {
	l := len(tag) + 2
	b := make([]byte, l)
	b[0] = '<'
	copy(b[1:l-1], tag)
	b[l-1] = '>'

	return b
}

func endTag(tag string) []byte {
	l := len(tag) + 3
	b := make([]byte, l)
	b[0] = '<'
	b[1] = '/'
	copy(b[2:l-1], tag)
	b[l-1] = '>'

	return b
}
