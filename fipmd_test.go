package fipmd

import (
	"log"
	"testing"
)

/*
func TestQuote1(t *testing.T) {
	input := `>quoted line1.
>another line 2
`
	el, idx := parseQuote([]byte(input))
	buf := new(bytes.Buffer)
	el.html(buf)
	if string(buf.Bytes()) != "<blockquote><p>quoted line1.<br> another line 2</p></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", string(buf.Bytes()), idx, len(input))
		t.FailNow()
	}
}
*/

func TestParse(t *testing.T) {
	input := `## Header #
text

1. >ordered 1
continue quote

2. ordered 2

>quote

* unordered 1
* unordered 2
	`

	out := Parse([]byte(input))
	if string(out) != `<h2>Header</h2><p>text</p><ol><li><blockquote><p>ordered 1 continue quote</p></blockquote></li><li>ordered 2</li></ol><blockquote><p>quote</p></blockquote><ul><li>unordered 1</li><li>unordered 2</li></ul>` {
		t.Log("TestParse, output:", string(out))
		t.Fail()
	}
}

func hook1(header *Header) {
	log.Println("hook1", header.Tag(), header.Text())
}

func hook2(header *Header) {
	log.Println("hook2", header.Tag(), header.Text())
}

func TestHook(t *testing.T) {
	input := `# Header1 #
	## header2`
	parser := NewParser()
	parser.AddHeaderHook(hook1)
	parser.AddHeaderHook(hook2)

	result := parser.Parse([]byte(input))
	t.Log(string(result))

	/*if string(buf.Bytes()) != "<h2>Header</h2>" ||
		idx != 12 {
		t.Log("TestHeader, output:",string(buf.Bytes()),idx)
		t.Fail()
	}*/
}

func TestA(t *testing.T) {
	input := `>**NOTE**
>The operator -> marks a pair, a -> b equals to (a, b)`
	/*2. A statement is considered not end when
	&nbsp;&nbsp;The line ends in the middle of parentheses() or brackets[]
	&nbsp;&nbsp;The line ends in a word that not valid as the end of a statement, such as infix operators (+, -, \*, / ...)
	*/
	//out := Parse([]byte(input))
	//t.Logf("%s",out)
	el, _ := parseQuote([]byte(input), nil)
	html := el.Html()
	t.Log(html)

}

func TestParseFile(t *testing.T) {
	path := "/work/fip/scala-basic/src/content/en/7.4.map.md"
	c, e := ParseFile(path)
	if e != nil {
		t.Log("err:", e)
	} else {
		t.Log(string(c))
	}

}
