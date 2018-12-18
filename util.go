package fipmd

import (
	"bytes"
	. "github.com/fipress/fiputil"
)

func htmlString(el markdown) string {
	buf := new(bytes.Buffer)
	el.html(buf)
	return buf.String()
}

var escapeChars = []byte("\\`*_{}[]()#+-.!:|&<>~")

func shouldEscape(input []byte) bool {
	if len(input) > 1 &&
		input[0] == '\\' &&
		bytes.IndexByte(escapeChars, input[1]) != -1 {
		return true
	}
	return false
}

func isIndentCode(input []byte) (isIndent bool, skip int) {
	isIndent = false
	if input[0] == '\t' {
		isIndent = true
		skip = 1 + skipChar(input[1:], '\t')
	} else if input[0] == ' ' {
		skip = 1 + skipChar(input[1:], ' ')
		if skip > 3 {
			isIndent = true
		}
	}
	return
}

func isSpan(c byte) bool {
	switch c {
	case '[', '*', '~', '!', '`', ' ', '<':
		return true
	}

	return false
}

func isLineEnd(input []byte) (bool, int) {
	switch input[0] {
	case '\n', '\r', '\f':
		return true, 1
	}
	return false, 0
}

func skipLeft(input []byte) (skip int) {
	i := 0
L:
	for ; i < len(input); i++ {
		switch input[i] {
		case ' ', '\t', '\r', '\n', '\f':
			continue L
		default:
			return i
		}
	}
	return i
}

func skipSpace(input []byte) (skip int) {
	i := 0
L:
	for ; i < len(input); i++ {
		switch input[i] {
		case ' ', '\t':
			continue L
		default:
			return i
		}
	}
	return i
}

func skipChar(input []byte, c byte) int {
	for i := 0; i < len(input); i++ {
		if input[i] != c {
			return i
		}
	}
	return 0
}

func skipUntilLineEnd(input []byte) (idx int, found bool) {
	i := 0
	for ; i < len(input); i++ {
		if input[i] == '\n' || input[i] == '\r' || input[i] == '\f' {
			return i, true
		}
	}
	return i, false
}

/*func skipUntilFunc(input []byte, func f(c byte) bool) (idx int,found bool) {

}*/

func skipUntilOrStopAt(input []byte, until byte, stop byte) (idx int, found bool) {
	i := 0
	for ; i < len(input); i++ {
		if input[i] == until {
			return i, true
		} else if input[i] == stop {
			return i, false
		}
	}
	return i, false
}

func skipUntilCharOrBlockEnd(input []byte, until byte) (idx int, found bool) {
	i := 0
	for i < len(input) {
		if input[i] == until {
			return i + 1, true
		} else if i+1 < len(input) && IsBlockEnd(input[i+1:]) {
			return i, false
		} else {
			i++
		}
	}
	return i, false
}

func SkipUntilArrayOrBlockEnd(input []byte, until []byte) (idx int, found bool) {
	i := 0
	l := len(until)
	for i < len(input) {
		//fmt.Printf("i=%d,input[i]=%c,until[0]=%c.\n",i,input[i],until[0])
		if i+l <= len(input) && bytes.Compare(input[i:i+l], until) == 0 {
			return i + l, true
		} else if i+i <= len(input) && IsBlockEnd(input[i+1:]) {
			return i, false
		} else {
			i++
		}
	}
	return i, false
}

func applySkipFuncs(input []byte, skipFuncs []func([]byte) int) (idx int) {
	for j := 0; j < len(skipFuncs); j++ {
		if skipFuncs[j] != nil {
			idx += skipFuncs[j](input[idx:])
		}
	}
	return idx
}

//extract text, treat '\n' as ' '(space)
func extractTextUntilFunc(input []byte, endFunc func([]byte) (bool, int), skipFuncs []func([]byte) int) (found bool, text []byte, idx int) {
	buf := new(bytes.Buffer)
	from := -1
	for idx < len(input) {
		shouldEnd, delta := endFunc(input[idx:])
		//fmt.Printf("extractTextUntilFunc,len=%d,i=%d,input=%c\n",len(input),idx,input[idx])
		if shouldEnd {
			if from != -1 {
				buf.Write(input[from:idx])
			}
			return true, buf.Bytes(), idx + delta
		}

		if input[idx] == '\n' {
			if from != -1 {
				buf.Write(input[from:idx])
				from = -1
			}

			if IsBlockEnd(input[idx:]) {
				return false, buf.Bytes(), idx + 1
			} else {
				buf.WriteByte(' ')
				idx++
				idx += applySkipFuncs(input[idx:], skipFuncs)
			}
		} else {
			if from == -1 {
				from = idx
			}
			idx++
		}

	}

	return false, nil, 0

}

func extractTextUntil(input []byte, c byte, skipFuncs []func([]byte) int) (bool, []byte, int) {
	endFunc := func(input []byte) (bool, int) {
		if input[0] == c {
			return true, 1
		} else {
			return false, 0
		}
	}

	return extractTextUntilFunc(input, endFunc, skipFuncs)
}

func extractTextUntilArray(input []byte, array []byte, skipFuncs []func([]byte) int) (bool, []byte, int) {
	endFunc := func(input []byte) (bool, int) {
		l := len(array)
		if bytes.Compare(input[:l], array) == 0 {
			return true, l
		} else {
			return false, 0
		}
	}

	return extractTextUntilFunc(input, endFunc, skipFuncs)
}

func extractTextUntilSpace(input []byte, skipFuncs []func([]byte) int) (bool, []byte, int) {
	endFunc := func(input []byte) (bool, int) {
		if IsLineEnd(input[0]) || IsSpace(input[0]) {
			return true, 1
		} else {
			return false, 0
		}
	}

	return extractTextUntilFunc(input, endFunc, skipFuncs)
}

/*
func extractUntilCharOrBlockEnd(input []byte, c byte,skipFuncs []func([]byte) int) (idx int, found bool) {
	i := 0
	for i<len(input) {
		if input[i] == until {
			return i+1, true
		} else if i+1 < len(input) && IsBlockEnd(input[i+1:]) {
			return i, false
		} else {
			i++
		}
	}
	return i, false

	endFunc := func(input []byte) (bool, int) {
		if input[0] == c {
			return true, 1
		} else if IsBlockEnd(input) {
			return false, 0
		} else {
			return false, 0
		}
	}

	return extractTextUntilFunc(input,endFunc,skipFuncs)
}*/
