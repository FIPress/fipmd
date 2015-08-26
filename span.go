package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil"
)

type span struct {
	tag     string
	content []byte
}

func (s *span) html(buf *bytes.Buffer) {
	if len(s.tag) != 0 {
		buf.Write(startTag(s.tag))
	}
	if len(s.content) != 0 {
		buf.Write(s.content)
	}
	if len(s.tag) != 0 {
		buf.Write(endTag(s.tag))
	}
}

func (s *span) Html() string {
	return htmlString(s)
}

func (s *span) inText() []byte {
	return s.content
}

func (s *span) Text() string {
	return string(s.content)
}

func (s *span) execHooks(parser *Parser) {
	//should implement by holders
}

type Plain struct {
	*span
}

func NewPlain(content []byte) *Plain {
	return &Plain{&span{"", content}}
}

func (p *Plain) execHooks(parser *Parser) {
	for _, hk := range parser.plainHooks {
		hk(p)
	}
}

type RawHtml struct {
	*span
}

func NewRawHtml(content []byte) *RawHtml {
	return &RawHtml{&span{"", content}}
}

func (r *RawHtml) execHooks(parser *Parser) {
	for _, hk := range parser.rawHtmlHooks {
		hk(r)
	}
}

type Link struct {
	href  []byte
	text  []byte
	title []byte
}

func (l *Link) html(buf *bytes.Buffer) {
	buf.WriteString(`<a href="`)
	//buf.Write(l.href)
	HtmlEscapeToBuffer(l.href, buf)
	buf.WriteByte('"')
	if len(l.title) != 0 {
		buf.WriteString(` title="`)
		//buf.Write(l.title)
		HtmlEscapeToBuffer(l.title, buf)
		buf.WriteByte('"')
	}
	buf.WriteByte('>')
	//buf.Write(l.text)
	HtmlEscapeToBuffer(l.text, buf)
	buf.Write(endTag("a"))
}

func (l *Link) Html() string {
	return htmlString(l)
}

func (l *Link) inText() []byte {
	return l.text
}

func (l *Link) Text() string {
	return string(l.text)
}

func (l *Link) execHooks(parser *Parser) {
	for _, hk := range parser.linkHooks {
		hk(l)
	}
}

type Img struct {
	*Link
}

func (i *Img) html(buf *bytes.Buffer) {
	buf.WriteString(`<img src="`)
	//buf.Write(i.href)
	HtmlEscapeToBuffer(i.href, buf)
	buf.WriteByte('"')
	if len(i.title) != 0 {
		buf.WriteString(` title="`)
		//buf.Write(i.title)
		HtmlEscapeToBuffer(i.title, buf)
		buf.WriteByte('"')
	}
	if len(i.text) != 0 {
		buf.WriteString(` alt="`)
		//buf.Write(i.text)
		HtmlEscapeToBuffer(i.text, buf)
		buf.WriteByte('"')
	}
	buf.WriteByte('>')
}

func (i *Img) Html() string {
	return htmlString(i)
}

func (i *Img) inText() []byte {
	return i.text
}

func (i *Img) Text() string {
	return string(i.text)
}

func (i *Img) execHooks(parser *Parser) {
	for _, hk := range parser.imgHooks {
		hk(i)
	}
}

type Em struct {
	*span
}

func NewEm(content []byte) *Em {
	return &Em{&span{"em", content}}
}

func (e *Em) execHooks(parser *Parser) {
	for _, hk := range parser.emHooks {
		hk(e)
	}
}

type Strong struct {
	*span
}

func NewStrong(content []byte) *Strong {
	return &Strong{&span{"strong", content}}
}

func (s *Strong) execHooks(parser *Parser) {
	for _, hk := range parser.strongHooks {
		hk(s)
	}
}

type Del struct {
	*span
}

func NewDel(content []byte) *Del {
	return &Del{&span{"del", content}}
}

func (d *Del) execHooks(parser *Parser) {
	for _, hk := range parser.delHooks {
		hk(d)
	}
}

type Code struct {
	*span
}

func NewCode(content []byte) *Code {
	return &Code{&span{"code", content}}
}

func (c *Code) html(buf *bytes.Buffer) {
	//todo:unescape
	writeSpanLiterally("code", c.content, buf)
}

func (c *Code) execHooks(parser *Parser) {
	for _, hk := range parser.codeHooks {
		hk(c)
	}
}

type self struct {
	tag string
}

func (s *self) html(buf *bytes.Buffer) {
	buf.WriteByte('<')
	buf.WriteString(s.tag)
	buf.WriteByte('>')
}

func (s *self) Html() string {
	return htmlString(s)
}

func (s *self) inText() []byte {
	return []byte{}
}

func (s *self) Text() string {
	return ""
}

func (s *self) execHooks(parser *Parser) {
	//should implement by the holders
}

type Br struct {
	*self
}

func NewBr() *Br {
	return &Br{&self{"br"}}
}

func (b *Br) execHooks(parser *Parser) {
	for _, hk := range parser.brHooks {
		hk(b)
	}
}

type Hr struct {
	*self
}

func NewHr() *Hr {
	return &Hr{&self{"hr"}}
}

func (h *Hr) execHooks(parser *Parser) {
	for _, hk := range parser.hrHooks {
		hk(h)
	}
}
