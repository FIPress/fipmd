package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil
)

/*func parsePlainUntil(input []byte, until byte,skipFuncs ...func([]byte) int) (ps []*plain, idx int) {
	ps = new([]*plain)
	from := -1

	for i := 0; i<len(input); {
		switch input[i] {
		case '\n':
			if from != -1 {
				input[i] = ' '
				ps = append(ps,&plain{input[from:i+1]})

			}
			i++
			i += applySkipFuncs(input[i:],skipFuncs)
		case until:
			if from != -1 {
				ps = append(ps,&plain{input[from:i]})
			}
			return ps,i
		default:
			if from == -1 {
				from = i
				i++
			}
		}
	}
	return nil, -1
}*/

func parseInlineLink(input []byte, skipFuncs []func([]byte) int) (l *Link, idx int) {
	if len(input) < 3 {
		return nil, 0
	}
	found, text, delta := extractTextUntil(input[1:], ']', skipFuncs)
	//fmt.Println("parseInlineLink:extractTextUntil,text,delta:",delta)
	if !found || delta == 1 {
		return nil, 0
	}
	idx += delta + 1
	if input[idx] == '(' {
		idx++
	} else {
		return nil, 0
	}
	//fmt.Printf("parseInlineLink,idx=%d,c[idx-1]=%c,c[idx]=%c,c[idx+1]=%c,\n",idx,input[idx-1],input[idx],input[idx+1])
	hrefFrom := idx
	hrefTo := 0
	/*if len(input) <= hrefFrom {
		return nil, 0
	}*/

	l = new(Link)
L:
	for i := hrefFrom; i < len(input); {
		switch input[i] {
		/*case '\n':
		input[i] = ' '*/
		case ' ', '\n':
			if IsBlockEnd(input[i:]) {
				return nil, 0
			}
			hrefTo = i
			delta, found := SkipUntil(input[i:], '"')
			if !found {
				continue L
			}
			i += delta + 1
			//fmt.Printf("delta=%d,i=%d,c=%c\n",delta,i,input[i])
			found, plains, delta := extractTextUntil(input[i:], '"', skipFuncs)
			//fmt.Printf("plains=%v,delta=%d,i=%d,c=%c\n",plains,delta,i,input[delta+i])
			if !found {
				return nil, 0
			}
			i += delta
			l.title = plains
		case ')':
			if hrefTo == 0 {
				hrefTo = i
			}
			idx = i + 1
			break L
		default:
			i++
		}
	}

	if hrefTo == 0 {
		return nil, 0
	}

	l.text = text
	l.href = input[hrefFrom:hrefTo]

	return l, idx
}

func parseEmphasis(input []byte, skipFuncs []func([]byte) int) (el markdown, idx int) {
	if len(input) > 4 && input[1] == '*' {
		el, idx = parseStrong(input[2:], skipFuncs)
		if idx != 0 {
			idx += 2
		}
	} else {
		el, idx = parseEm(input[1:], skipFuncs)
		if idx != 0 {
			idx++
		}
	}
	//fmt.Println("parseEmphasis:el,idx",el, idx)

	return
}

func parseEm(input []byte, skipFuncs []func([]byte) int) (e *Em, idx int) {
	//fmt.Printf("parseEm before, input=%s,len(f)=%d\n",input,len(skipFuncs))
	found, content, idx := extractTextUntil(input, '*', skipFuncs)
	//fmt.Println("parseEm, after:c=",content,idx)
	if !found || idx == 0 || len(content) == 0 {
		return nil, 0
	} else {
		return NewEm(content), idx
	}
}

func parseStrong(input []byte, skipFuncs []func([]byte) int) (*Strong, int) {
	found, content, idx := extractTextUntilArray(input, []byte{'*', '*'}, skipFuncs)
	if !found || idx == 0 || len(content) == 0 {
		return nil, 0
	} else {
		return NewStrong(content), idx
	}
}

func parseBr(input []byte) (*Br, int) {
	if len(input) < 3 {
		return nil, 0
	}
	if bytes.Compare(input[:3], []byte{' ', ' ', '\n'}) == 0 {
		return NewBr(), 2
	} else {
		return nil, 0
	}
}

func parseHr(input []byte) (*Hr, int) {
	l := len(input)
	var c byte
	count := 0
	i := 1

	if l < 3 {
		goto NotHr
	}

	c = input[0]
	if !isValidHr(c) {
		goto NotHr
	}

	for ; i < l; i++ {
		switch input[i] {
		case c:
			count++
		case ' ':
		case '\n':
			goto Hr
		default:
			goto NotHr
		}
	}
Hr:
	if count > 1 {
		return NewHr(), i
	} else {
		return nil, 0
	}
NotHr:
	//input[0] = ' '
	return nil, 0

}

func parseDel(input []byte, skipFuncs []func([]byte) int) (markdown, int) {
	found, content, idx := extractTextUntilArray(input, []byte{'~', '~'}, skipFuncs)
	if !found || idx == 0 {
		return nil, 0
	} else {
		return NewDel(content), idx
	}
}

func parseImage(input []byte, skipFuncs []func([]byte) int) (markdown, int) {
	a, idx := parseInlineLink(input[1:], skipFuncs)
	if a == nil {
		return nil, 0
	}

	idx++
	img := &Img{a}

	return img, idx
}

func parseInlineCode(input []byte, skipFuncs []func([]byte) int) (markdown, int) {
	idx, found := SkipUntil(input[1:], '`')
	if !found {
		return nil, 0
	}

	idx++
	e := NewCode(input[1:idx])
	idx++
	return e, idx
}
