package fipmd

import (
	"bytes"
	"strconv"
)

type container interface {
	add(el markdown)
}

type block struct {
	tag string
	els []markdown
}

func newBlock(tag string) *block {
	return &block{tag, make([]markdown, 0)}
}

func (b *block) html(buf *bytes.Buffer) {
	buf.Write(startTag(b.tag))
	for _, el := range b.els {
		if el != nil {
			el.html(buf)
		}
	}
	buf.Write(endTag(b.tag))
}

func (b *block) Html() string {
	return htmlString(b)
}

func (b *block) inText() []byte {
	buf := new(bytes.Buffer)
	for _, el := range b.els {
		buf.Write(el.inText())
	}
	return buf.Bytes()
}

func (b *block) Text() string {
	return string(b.inText())
}

func (b *block) execHooks(parser *Parser) {
	//not implemented
}

func (b *block) add(el markdown) {
	/*l := len(b.els)
	if l > 1 {
		last := b.els[l-1]
		switch lb := last.(type){
		case *block:
			switch b:=el.(type) {
			case *block:
				if b.tag == "p" && lb.tag == "p" {
					lb.addAll(b.els)
				}
			}
		}


	}*/

	b.els = append(b.els, el)
}

func (b *block) addAll(el []markdown) {
	b.els = append(b.els, el...)
}

type Para struct {
	*block
}

func NewPara() *Para {
	return &Para{newBlock("p")}
}

func (p *Para) execHooks(parser *Parser) {
	for _, hk := range parser.paraHooks {
		hk(p)
	}
}

type Header struct {
	*block
	Level int
}

func NewHeader(level int) *Header {
	return &Header{newBlock("h" + strconv.Itoa(level)), level}
}

func (h *Header) Tag() string {
	return h.block.tag
}

func (h *Header) execHooks(parser *Parser) {
	for _, hk := range parser.headerHooks {
		hk(h)
	}
}

type Quote struct {
	*block
}

func NewQuote() *Quote {
	return &Quote{newBlock("blockquote")}
}

func (q *Quote) execHooks(parser *Parser) {
	for _, hk := range parser.quoteHooks {
		hk(q)
	}
}

type Li struct {
	*block
}

func NewLi() *Li {
	return &Li{newBlock("li")}
}

func (l *Li) execHooks(parser *Parser) {
	for _, hk := range parser.liHooks {
		hk(l)
	}
}

type Ol struct {
	*block
}

func NewOl() *Ol {
	return &Ol{newBlock("ol")}
}

func (o *Ol) execHooks(parser *Parser) {
	for _, hk := range parser.olHooks {
		hk(o)
	}
}

type Ul struct {
	*block
}

func NewUl() *Ul {
	return &Ul{newBlock("ul")}
}

func (u *Ul) execHooks(parser *Parser) {
	for _, hk := range parser.ulHooks {
		hk(u)
	}
}

type CodeBlock struct {
	*block
}

func NewCodeBlock(code markdown) *CodeBlock {
	//b := &block{"pre",make([]markdown, 0)}
	//b.add(code)
	return &CodeBlock{&block{"pre", []markdown{code}}}
}

func (cb *CodeBlock) execHooks(parser *Parser) {
	for _, hk := range parser.codeBlockHooks {
		hk(cb)
	}
}
