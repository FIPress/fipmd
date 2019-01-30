package fipmd

import (
	"bytes"
	"io/ioutil"
)

/*var htmlTags = []string{
"p":          true,
"dl":         true,
"h1":         true,
"h2":         true,
"h3":         true,
"h4":         true,
"h5":         true,
"h6":         true,
"ol":         true,
"ul":         true,
"del":        true,
"div":        true,
"ins":        true,
"pre":        true,
"form":       true,
"math":       true,
"table":      true,
"iframe":     true,
"script":     true,
"fieldset":   true,
"noscript":   true,
"blockquote": true,

// HTML5
"video":      true,
"aside":      true,
"canvas":     true,
"figure":     true,
"footer":     true,
"header":     true,
"hgroup":     true,
"output":     true,
"article":    true,
"section":    true,
"progress":   true,
"figcaption": true,
}}*/

func ParseFile(path string) (content []byte, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	content = Parse(bytes)

	return
}

func ParseString(input string) string {
	return string(Parse([]byte(input)))
}

func Parse(input []byte) []byte {
	els := parse(input)

	if len(els) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	for _, el := range els {
		//fmt.Println("el",el)
		el.html(buf)
	}
	return buf.Bytes()
}

type Parser struct {
	headerHooks    []func(*Header)
	paraHooks      []func(*Para)
	quoteHooks     []func(*Quote)
	ulHooks        []func(*Ul)
	olHooks        []func(*Ol)
	liHooks        []func(*Li)
	codeBlockHooks []func(*CodeBlock)
	plainHooks     []func(*Plain)
	rawHtmlHooks   []func(*RawHtml)
	linkHooks      []func(*Link)
	imgHooks       []func(*Img)
	codeHooks      []func(*Code)
	emHooks        []func(*Em)
	strongHooks    []func(*Strong)
	delHooks       []func(*Del)
	brHooks        []func(*Br)
	hrHooks        []func(*Hr)
}

func NewParser() *Parser {
	return new(Parser)
}

func (parser *Parser) AddHeaderHook(headerHook func(*Header)) {
	parser.headerHooks = append(parser.headerHooks, headerHook)
}

func (parser *Parser) AddParaHook(hook func(*Para)) {
	parser.paraHooks = append(parser.paraHooks, hook)
}

func (parser *Parser) Parse(input []byte) []byte {
	els := parse(input)

	if len(els) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	for _, el := range els {
		el.execHooks(parser)
		el.html(buf)
	}
	return buf.Bytes()
}
