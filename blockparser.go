package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil"
)

const (
	indentSize = 4
)

var (
	fenced = []byte{'`', '`', '`'}
)

func parseRawHtml(input []byte, tag []byte, selfClosing bool) (*RawHtml, int) {
	i := 1 + len(tag)
	delta, found := 0, false
	if selfClosing {
		delta, found = skipUntilCharOrBlockEnd(input[i:], '>')
	} else {
		buf := new(bytes.Buffer)
		buf.Write([]byte{'<', '/'})
		buf.Write(tag)
		buf.WriteRune('>')
		close := buf.Bytes()
		//delta,found = SkipUntilArrayOrBlockEnd(input[i:],close)

		delta, found = SkipUntilArray(input[i:], close)
		delta += len(close)

		//fmt.Println("delta:",delta,"found:",found,close,input[i:])
	}
	if found {
		i += delta
		return NewRawHtml(input[:i]), i
	} else {
		return nil, 0
	}
}

func parseIndentCode(input []byte) markdown {
	el := new(Plain) //todo

	return el
}

func parsePara(input []byte, isBlock bool, skipFuncs ...func([]byte) int) (*Para, int) {
	p := NewPara()
	els, idx := parseContent(input, isBlock, skipFuncs...)
	p.addAll(els)
	return p, idx
}

/*func parseLine(input []byte,skipFuncs ...func([]byte) int) ( *block, int) {
	p := NewPara()
	els,idx := parseContent(input,false,skipFuncs...)
	p.addAll(els)
	return p,idx
}*/

func parseCode(input []byte) (markdown, int) {
	if len(input) > 6 {
		//buf := new(bytes.Buffer)
		i := 3
		lang := ""
		if IsLetter(input[3]) {
			found, langBytes, delta := extractTextUntilSpace(input[3:], nil)
			if found && len(langBytes) != 0 {
				lang = string(langBytes)
				i += delta
			}
		}

		i += skipLeft(input[i:])

		delta, found := SkipUntilArray(input[i:], []byte{'\n', '`', '`', '`'})
		if found {
			from := i
			end := from + delta
			i += delta + 4
			empty, delta := IsBlankLine(input[i:])
			if empty {
				if lang == "" {
					return NewPre([]markdown{NewCode(input[from:end])}), i + delta
				} else {
					return NewPre([]markdown{NewCodeWithLang(input[from:end], lang)}), i + delta
				}
			}

		}

		/*for i < len(input) {
			//i += applySkipFuncs(input[i:], skipFuncs)
			if prefixFenced(input[i:]) {
				i += 3
				i += SkipSpaceAndLineEnd(input[i:])
				if lang == "" {
					return NewPre([]markdown{NewCode(buf.Bytes())}), i
				} else {
					return NewPre([]markdown{NewCodeWithLang(buf.Bytes(), lang)}), i
				}
			} else if buf.Len() != 0 {
				buf.WriteByte('\n')
			}
			delta, found := skipUntilLineEnd(input[i:])
			//fmt.Printf("parseCode:i=%d,c=%c,delta=%d.\n",i,input[i],delta)
			if found {
				if delta != 0 {
					from := i
					i += delta
					buf.Write(input[from:i])
				}
				i++
			} else {
				return nil, 0
			}
		}*/
	}

	return nil, 0
}

/*func skipUntilFencedCodeLineEnd(input []byte) (idx int, found bool) {
	i:=0
	for ;i<len(input);i++ {
		if input[i] == '\n' || input[i] == '\r' || input[i] == '\f'{
			return i+1, true
		} else if fiputil.SliceEquals(input[i:i+3],fenced) {
			return i, true
		}
	}
	return 0, false
}*/

func parseHeader(input []byte) (el *Header, idx int) {
	level := skipChar(input, '#')
	if level == 0 {
		return
	}
	to := level + skipUntilHeaderClosing(input[level:])
	els, _ := parseContent(input[level:to], false)

	el = NewHeader(level)
	el.addAll(els)

	idx = to + skipHeaderClosing(input[to:])

	return
}

func skipUntilHeaderClosing(input []byte) int {
	i := 0
	l := len(input)
	for ; i < l; i++ {
		switch input[i] {
		case ' ':
			if i+1 < l && input[i+1] == '#' {
				return i
			}
		case '#':
			return i
		case '\n', '\r', '\f':
			return i
		}
	}
	return i

}

func skipHeaderClosing(input []byte) int {
	i := 0
	for i < len(input) {
		switch input[i] {
		case ' ', '\t', '#':
			i++
		case '\n', '\r', '\f':
			return i + 1
		default:
			return 0
		}
	}
	return i
}

func parseQuote(input []byte, skipFuncs ...func([]byte) int) (el *Quote, idx int) {
	el = NewQuote()

	for idx < len(input) {
		//fmt.Printf("parseQuote:before skip,idx=%d,c=%c.\n",idx,input[idx])
		idx += skipQuotePrefix(input[idx:])
		//fmt.Println("parseQuote:before add func",skipFuncs,skipQuotePrefix)
		sub, delta := parseBlock(input[idx:], true, append(skipFuncs, skipQuotePrefix)...)
		if delta == 0 {
			return
		}

		el.add(sub)
		idx += delta
		if idx >= len(input) {
			return
		}
		//fmt.Printf("parseQuote:before should end ,idx=%d,c=%c.\n",idx,input[idx])
		if shouldQuoteEnd(input[idx:], skipFuncs) {
			return
		}
	}

	return
}

func skipQuotePrefix(input []byte) int {
	i := 0
	count := 0
	for i < len(input) {
		switch input[i] {
		case '>':
			if count == 0 {
				count++
				i++
			} else {
				return i
			}
		case ' ', '\t':
			i++
		default:
			return i
		}
	}
	return i
}

func shouldQuoteEnd(input []byte, skipFuncs []func([]byte) int) bool {
	idx := applySkipFuncs(input, skipFuncs)
	idx += skipLeft(input[idx:])
	if len(input[idx:]) == 0 {
		return true
	}
	return input[idx] != '>'
}

//todo: 1. If list items are separated by blank lines, Markdown will wrap the items in <p>tags in the HTML outputList
//  2. items may consist of multiple paragraphs.
// Each subsequent paragraph in a list item must be indented by either 4 spaces or one tab:
func parseLi(input []byte, valid func([]byte) bool, skipFuncs []func([]byte) int) (el *Li, idx int) {
	el = NewLi()
	for idx < len(input) {
		sub, delta := parseBlock(input[idx:], false, skipFuncs...)
		//fmt.Printf("parseLi:after parse block ,sub=%v,delta=%d.idx=%d,len(input)=%d\n",sub.Html(),delta,idx,len(input))
		idx += delta
		if idx > len(input) {
			return
		}
		shouldEnd, isPara, delta, addSkip := shouldLiEnd(input[idx:], valid, skipFuncs)
		if addSkip {
			skipFuncs = append(skipFuncs, skipIndent)
		}
		//fmt.Printf("parseLi:shouldend=%v,ispara=%v,delta=%d,add=%v.\n",shouldEnd,isPara,delta,addSkip)
		if isPara {
			el.add(sub)
		} else {
			//p := *para(sub)
			switch b := sub.(type) {
			case *Para:
				el.addAll(b.els)
			default:
				el.add(b)
			}
		}

		idx += delta
		//fmt.Printf("parseQuote:before should end ,idx=%d,c=%c.\n",idx,input[idx])
		if shouldEnd {
			return
		}
	}
	return
}

func parseUl(input []byte, skipFuncs []func([]byte) int) (el *Ul, idx int) {
	el = NewUl()
	idx = parseList(input, el, isUl, skipFuncs)
	return
}

func parseOl(input []byte, skipFuncs []func([]byte) int) (el *Ol, idx int) {
	el = NewOl()
	idx = parseList(input, el, isOl, skipFuncs)
	return
}

func parseList(input []byte, parent container, valid func([]byte) bool, skipFuncs []func([]byte) int) (idx int) {
	i := 0
	for i < len(input) {
		i += skipLeft(input[i:])
		if valid(input[i:]) {
			i += 2
			i += skipSpace(input[i:])
			li, delta := parseLi(input[i:], valid, skipFuncs)
			i += delta
			parent.add(li)

		} else {
			break
		}

		end, delta := shouldListEnd(input[i:], valid, skipFuncs)
		if end {
			return i + delta
		}
	}
	return i
}

func skipIndent(input []byte) int {
	if len(input) == 0 {
		return 0
	}
	switch input[0] {
	case '\t':
		return 1
	case ' ':
		if IsSpaceIndented(input, indentSize) {
			return indentSize
		}
	}
	return 0
}

//1. end of input sure means end of li
//2. if there is a blank line, then
//	a. next line is indented, means it consists multiple paragraphs, so, it's not end yet
//	b. next line is not indented, end
//3. indented line indicate that there might be nested item
func shouldLiEnd(input []byte, valid func([]byte) bool, skipFuncs []func([]byte) int) (shouldEnd bool, isPara bool, idx int, addSkipFunc bool) {
	//fmt.Printf("shouldLiEnd, invoke apply, len of skiper: %d, input=%s\n",len(skipFuncs),input[idx:])
	idx = applySkipFuncs(input, skipFuncs)
	if len(input[idx:]) == 0 {
		return true, false, 0, false
	}

	isBlank, delta := IsBlankLine(input[idx:])
	if isBlank {
		idx += delta
		if IsIndented(input[idx:], indentSize) {
			isPara = true
			shouldEnd = false
		} else {
			isPara = false
			shouldEnd = true
		}
	} else {
		isPara = false
		if IsIndented(input[idx:], indentSize) {
			//&& valid(input[idx:])
			addSkipFunc = true
			shouldEnd = false
		} else {
			shouldEnd = true
		}
	}

	return
}

func shouldListEnd(input []byte, valid func([]byte) bool, skipFuncs []func([]byte) int) (bool, int) {
	idx := applySkipFuncs(input, skipFuncs)
	idx += skipLeft(input[idx:])
	return !valid(input[idx:]), idx
}

func isUl(input []byte) bool {
	if len(input) > 2 &&
		isBullet(input[0]) &&
		input[1] == ' ' {
		return true
	}

	return false
}

func isBullet(c byte) bool {
	switch c {
	case '*', '+', '-':
		return true
	}
	return false
}

func isOl(input []byte) bool {
	if len(input) > 2 &&
		isNumber(input[0]) &&
		input[1] == '.' {
		return true
	}

	return false
}

func isNumber(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}

	return false
}
