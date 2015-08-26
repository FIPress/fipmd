package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil"
)

/*const (
	indentSize = 4
)

var (
	//inlineLinkR = regexp.MustCompile(`\[(.*)\]\((\S*)\s+"(.*)"\)`)

)*/

func parse(input []byte) (els []markdown) {
	els = make([]markdown, 0, 5)
	for i := 0; i < len(input); {
		el, delta := parseBlock(input[i:], false)
		if el != nil {
			els = append(els, el)
			i += delta
		} else {
			i++
		}
	}
	return
}

func parseBlock(input []byte, isBlock bool, skipFuncs ...func([]byte) int) (markdown, int) {
	if len(input) == 0 {
		return nil, 0
	}

	idx := 0

	//fmt.Printf("parseBlock: %s\n",input)
	/*if prefixIndent(input) {
		//parse indent code
	}*/

	idx += SkipSpaceAndLineEnd(input)
	if idx >= len(input) {
		return nil, 0
	}
	//fmt.Printf("parseBlock:idx=%d,c=%c\n",idx,input[idx])
	tag, selfClosing := GetHtmlTag(input)
	//fmt.Println("parseBlock:GetHtmlTag:",tag)
	if len(tag) != 0 {
		el, delta := parseRawHtml(input, tag, selfClosing)
		//fmt.Println("parseBlock:parseRawHtml,el:",el,"delta:",delta)
		if delta != 0 {
			return el, idx + delta
		}
	}
	//fmt.Printf("parseBlock prefixFenced:idx=%d,c=%c\n",idx,input[idx])
	if prefixFenced(input[idx:]) {
		el, delta := parseCode(input[idx:])
		//		fmt.Println("parseBlock:prefixFenced:el=",el)
		if delta != 0 {
			return el, idx + delta
		}
	}
	//fmt.Printf("parseBlock prefixFenced:idx=%d,c=%c\n",idx,input[idx])
	if prefixHr(input[idx:]) {
		el, delta := parseHr(input[idx:])
		//		fmt.Println("parseBlock:prefixHr:el=",el,delta)
		if delta != 0 {
			return el, idx + delta
		}
	}

	//fmt.Printf("parseBlock prefixHr:idx=%d,c=%c\n",idx,input[idx])

	if prefixHeader(input[idx:]) {
		el, delta := parseHeader(input[idx:])
		//		fmt.Println("parseBlock:parseHeader:el=",el)
		if delta != 0 {
			return el, idx + delta
		}
	}

	if prefixQuote(input[idx:]) {
		el, delta := parseQuote(input[idx:], skipFuncs...)
		//		fmt.Println("parseBlock:parseQuote:el=",el)
		if delta != 0 {
			return el, idx + delta
		}
	}
	//fmt.Printf("parseBlock prefixUl:idx=%d,c=%c\n",idx,input[idx])
	if prefixUl(input[idx:]) {
		//fmt.Println("parseBlock:parseUl:idx=",idx)
		el, delta := parseUl(input[idx:], skipFuncs)
		//fmt.Println("parseBlock:parseUl:el=",el)
		if delta != 0 {
			return el, idx + delta
		}
	}
	//fmt.Printf("parseBlock prefixOl:idx=%d,c=%c\n",idx,input[idx])
	if prefixOl(input[idx:]) {
		//fmt.Println("parseBlock:parseOl:idx=",idx)
		el, delta := parseOl(input[idx:], skipFuncs)
		//fmt.Println("parseBlock:parseOl:el=",el)
		if delta != 0 {
			return el, idx + delta
		}
	}
	//fmt.Println("parseBlock, parsePara, idx:",idx)
	//if not above, it is a normal paragraph
	/*if isLi {
		el ,delta := parseLine(input[idx:],skipFuncs...)
		fmt.Println("parseLine:parseContent:el=",el)
		return el, idx+delta
	} else {*/
	//fmt.Println("parseBlock,parsePara start")
	//fmt.Printf("parseBlock parsePara:idx=%d,c=%c\n",idx,input[idx])
	el, delta := parsePara(input[idx:], isBlock, skipFuncs...)
	//	fmt.Println("parseBlock:parsePara:el=",el)
	//fmt.Println("parseBlock,parsePara end,el=",el,"delta=",delta)
	return el, idx + delta
	//}
}

func parseContent(input []byte, isBlock bool, skipFuncs ...func([]byte) int) ([]markdown, int) {
	els := make([]markdown, 0)
	plainFrom := -1
	i := skipLeft(input)

	for i < len(input) {

		//	fmt.Printf("parseContent,i=%d,c=%c\n",i,input[i])
		if shouldEscape(input[i:]) {
			if plainFrom != -1 {
				els = append(els, NewPlain(input[plainFrom:i]))
			}
			plainFrom = i + 1
			i += 2
		} else if isSpan(input[i]) {
			el, delta := parseSpan(input[i:], skipFuncs)
			//fmt.Printf("parseContent,span,el=%v,delta=%d,input=%s\n",el,delta,input[i:])
			if delta != 0 {
				if plainFrom != -1 {
					els = append(els, NewPlain(input[plainFrom:i]))
					plainFrom = -1
				}
				els = append(els, el)
				i += delta
				//todo:i-1?
			} else {
				//fmt.Printf("parseContent,plainFrom=%d,c=%c\n",plainFrom,input[i])
				if plainFrom == -1 {
					plainFrom = i
				}
				i++
			}
		} else if IsLineEnd(input[i]) {
			if plainFrom != -1 {
				els = append(els, NewPlain(input[plainFrom:i]))
				plainFrom = -1
			}
			i++
			if i < len(input) {
				i += applySkipFuncs(input[i:], skipFuncs)
				if isBlock {
					shouldEnd, delta := IsBlankLine(input[i:])
					if shouldEnd {
						return els, i + delta
					}
					els = append(els, NewPlain([]byte{' '}))
				} else {
					//parse line by line
					return els, i
				}
			}
		} else {
			if plainFrom == -1 {
				plainFrom = i
			}
			i++
		}
	}

	if plainFrom != -1 {
		els = append(els, NewPlain(input[plainFrom:i]))
	}

	return els, i
}

func parseSpan(input []byte, skipFuncs []func([]byte) int) (el markdown, idx int) {
	switch input[0] {
	case '[':
		el, idx = parseInlineLink(input, skipFuncs)
	case '*':
		el, idx = parseEmphasis(input, skipFuncs)
	case '~':
		el, idx = parseDel(input, skipFuncs)
	case '!':
		el, idx = parseImage(input, skipFuncs)
	case '`':
		el, idx = parseInlineCode(input, skipFuncs)
	case ' ':
		el, idx = parseBr(input)
	case '<':
		tag, selfClosing := GetHtmlTag(input)
		if len(tag) != 0 {
			el, idx = parseRawHtml(input, tag, selfClosing)
		}
	case '\n':
		el, idx = NewPlain([]byte{' '}), 1

	}
	return
}

func prefixIndent(input []byte) bool {
	if input[0] == '\t' {
		return true
	} else if input[0] == ' ' && input[1] == ' ' && input[2] == ' ' && input[3] == ' ' {
		return true
	} else {
		return false
	}
}

func prefixFenced(input []byte) bool {
	if len(input) > 2 && bytes.Compare(input[:3], fenced) == 0 {
		return true
	} else {
		return false
	}
}

func prefixHeader(input []byte) bool {
	if input[0] == '#' {
		return true
	} else {
		return false
	}
}

func prefixHr(input []byte) bool {
	switch input[0] {
	case '-', '*', '_':
		return true
	}
	return false
}

func prefixQuote(input []byte) bool {
	if input[0] == '>' {
		return true
	} else {
		return false
	}
}

func prefixUl(input []byte) bool {
	if len(input) < 2 {
		return false
	}

	switch input[0] {
	case '*', '-', '+':
		if IsSpace(input[1]) {
			return true
		}
	}
	return false
}

func prefixOl(input []byte) bool {
	if len(input) < 2 {
		return false
	}

	if IsDigit(input[0]) &&
		input[1] == '.' {
		return true
	} else {
		return false
	}
}

func isValidHr(c byte) bool {
	switch c {
	case '-', '*', '_':
		return true
	}
	return false
}

//func isHr(input []byte) (bool)

/*func skipHr(input []byte) int {
	l := len(input)
	var c byte
	var i, count int

	if l < 4 {
		return 0
	}

	i = skipChar(input,' ')
	c = input[i]

	if !isValidHr(c) {
		return 0
	}

L:	for i++;i<l;i++ {
		switch input[i]{
		case c:
			count++
		case ' ':
		case '\n':
			i += 1
			break L
		default:
			return 0
		}
	}
//fmt.Println("count:",count)
	if count > 1 {
		return i
	} else {
		return 0
	}
}*/

//skip spaces and tabs
