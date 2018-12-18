package fipmd

import (
	"bytes"
	"strconv"
)

type container interface {
	add(el markdown)
}

type block struct {
	tag   string
	class string
	els   []markdown
}

func newBlock(tag string) *block {
	return &block{tag: tag, els: make([]markdown, 0)}
}

func newBlockWithClass(tag, class string) *block {
	return &block{tag, class, make([]markdown, 0)}
}

func (b *block) html(buf *bytes.Buffer) {
	if b.class == "" {
		buf.Write(startTag(b.tag))
	} else {
		buf.Write(startTagWithClass(b.tag, b.class))
	}

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

type Pre struct {
	*block
}

func NewPre(els []markdown) *Pre {
	return &Pre{&block{tag: "pre", els: els}}
}

type CodeBlock struct {
	*block
}

func NewCodeBlock(code markdown, lang string) *CodeBlock {
	//b := &block{"pre",make([]markdown, 0)}
	//b.add(code)
	if lang == "" {
		return &CodeBlock{&block{tag: "code", els: []markdown{code}}}
	} else {
		return &CodeBlock{&block{tag: "code", class: "language-" + lang, els: []markdown{code}}}
	}
}

func (cb *CodeBlock) execHooks(parser *Parser) {
	for _, hk := range parser.codeBlockHooks {
		hk(cb)
	}
}
