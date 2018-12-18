package fipmd

import (
	"testing"
)

func TestHtml(t *testing.T) {
	input := `<span>abcd</span>`
	el, idx := parseBlock([]byte(input), false)
	if el == nil {
		t.Log("parseHr failed")
		t.FailNow()
	}

	html := el.Html()
	if html != "<span>abcd</span>" ||
		idx != len(input) {
		t.Log("parseHtml failed, output:", html, idx)
		t.FailNow()
	}

	//fmt.Println("TestHtml 2")
	input = `<pre><code>println("abc")

	a = b+c
	</code></pre>`
	el, idx = parseBlock([]byte(input), false)
	//fmt.Println("TestHtml 3:el=",el,"idx=",idx)
	html = el.Html()
	if html != `<pre><code>println("abc")

	a = b+c
	</code></pre>` ||
		idx != len(input) {
		t.Log("parseHtml should ok, output:", html, idx)
		t.FailNow()
	}

	input = `<span>abcd</span`
	el, idx = parseBlock([]byte(input), false)
	//fmt.Println("TestHtml 3:el=",el,"idx=",idx)
	html = el.Html()
	if html != `<p><span>abcd</span</p>` ||
		idx != len(input) {
		t.Log("parseHtml should faile, and we should get an para instead, output:", html, idx)
		t.FailNow()
	}

	input = `<img src="abc">`
	el, idx = parseBlock([]byte(input), false)

	html = el.Html()
	if html != `<img src="abc">` ||
		idx != len(input) {
		t.Log("parseHtml failed, output:", html, idx)
		t.FailNow()
	}

	input = `<img src="abc"

>`
	el, idx = parseBlock([]byte(input), false)
	html = el.Html()
	if html != `<p><img src="abc"</p>` ||
		idx != len(input)-2 {
		t.Log("parseHtml should faile, and we should get an para instead, output:", html, idx, len(input))
		t.FailNow()
	}
}

func TestHr(t *testing.T) {
	input := `- - -`
	el, idx := parseHr([]byte(input))
	if el == nil {
		t.Log("parseHr failed")
		t.FailNow()
	}

	html := el.Html()
	if html != "<hr>" ||
		idx != len(input) {
		t.Log("parseHr failed, output:", html, idx)
		t.Fail()
	}
}

func TestHeader(t *testing.T) {
	input := `## Header #
	header2`
	el, idx := parseHeader([]byte(input))
	html := el.Html()
	if html != "<h2>Header</h2>" ||
		idx != 12 {
		t.Log("TestHeader, output:", html, idx)
		t.Fail()
	}
}

func TestLink(t *testing.T) {
	input := `[test](http://abc "title")as`
	el, idx := parseInlineLink([]byte(input), nil)
	html := el.Html()
	if html != `<a href="http://abc" title="title">test</a>` ||
		idx != len(input)-2 {
		t.Log("parseInlineLink, output:", html, idx, len(input))
		t.Fail()
	}

	input = `[test](http://abc title)as`
	el, idx = parseInlineLink([]byte(input), nil)
	html = el.Html()
	if html != `<a href="http://abc">test</a>` ||
		idx != len(input)-2 {
		t.Log("parseInlineLink, output:", html, idx, len(input))
		t.Fail()
	}

	input = `[test] (http://abc "title")as`
	el, idx = parseInlineLink([]byte(input), nil)
	if el != nil || idx != 0 {
		t.Error("Should fail, not a link. idx:", idx)
	}
	/*html = el.Html()
	if html != `<a href="http://abc" title="title">test</a>` ||
		idx != len(input) - 2 {
		t.Log("parseInlineLink, output:",html,idx,len(input))
		t.Fail()
	}*/
}

func TestImage(t *testing.T) {
	input := `![A link](http://address "title") ah`
	el, idx := parseImage([]byte(input), nil)
	html := el.Html()
	if html != `<img src="http://address" title="title" alt="A link">` ||
		idx != len(input)-3 {
		t.Log("parseImage, output:", html, idx)
		t.Fail()
	}

}

func TestPara(t *testing.T) {
	input := `para line1.
another line 2

e`
	el, idx := parsePara([]byte(input), true)
	html := el.Html()
	if html != "<p>para line1. another line 2</p>" ||
		idx != len(input)-1 {
		t.Log("parsePara, output:", html, idx, len(input))
		t.Fail()
	}
}

func TestCode(t *testing.T) {
	input := "```\nfunc Test(a int) {\n &copy;\nprint(a)}\n```"
	el, idx := parseCode([]byte(input))
	html := el.Html()
	if html != "<pre><code>func Test(a int) {\n &copy;\nprint(a)}</code></pre>" ||
		idx != len(input) {
		t.Log("parseCode failed, output:", html, idx, len(input))
		t.Fail()
	}
}

func TestCodeBlock1(t *testing.T) {
	input := "```" + `go
		println("a || b", a || b)
		println("a && b", a && b)
` + "```"
	el, idx := parseCode([]byte(input))
	html := el.Html()
	t.Log(idx)
	t.Log(html)
	/*if html != "<pre><code class=\"language-go\">func Test(a int) {\n &amp;copy;\nprint(a)}</code></pre>" ||
		idx != len(input) {
		t.Log("parseCode failed, output:", html, idx, len(input))
		t.Fail()
	}*/
}

func TestCodeWithLang(t *testing.T) {
	input := "```go func Test(a int) {\n &copy;\nprint(\"a\")}\n```"
	el, idx := parseCode([]byte(input))
	html := el.Html()

	if html != "<pre><code class=\"language-go\">func Test(a int) {\n &copy;\nprint(\"a\")}</code></pre>" ||
		idx != len(input) {
		t.Log("parseCode failed, output:", html, idx, len(input))
		t.Fail()
	}
}

func TestInlineCode(t *testing.T) {
	input := "`func Test(a int)`ab"
	el, idx := parseInlineCode([]byte(input), nil)
	html := el.Html()
	if html != "<code>func Test(a int)</code>" ||
		idx != len(input)-2 {
		t.Log("parseInlineCode failed, output:", html, idx, len(input))
	}
}

func TestQuote(t *testing.T) {
	input := `>quoted line1.
>another line 2
`
	el, idx := parseQuote([]byte(input))
	html := el.Html()
	if html != "<blockquote><p>quoted line1. another line 2</p></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", html, idx, len(input))
		t.FailNow()
	}

	input = `>quoted line1.
>another line 2
>
>>nested
>still nested
>`
	el, idx = parseQuote([]byte(input))
	html = el.Html()
	if html != "<blockquote><p>quoted line1. another line 2</p><blockquote><p>nested still nested</p></blockquote></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", html, idx)
		t.FailNow()
	}
	input = `>quoted **line1**.
>another line 2
>
>>nested
>
>return1`
	el, idx = parseQuote([]byte(input))
	html = el.Html()
	if html != "<blockquote><p>quoted <strong>line1</strong>. another line 2</p><blockquote><p>nested</p></blockquote><p>return1</p></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", html, idx)
		t.FailNow()
	}

	input = `>quoted line1.
>another line 2
>
>>nested
>still nested
>
>return1`
	el, idx = parseQuote([]byte(input))
	html = el.Html()
	if html != "<blockquote><p>quoted line1. another line 2</p><blockquote><p>nested still nested</p></blockquote><p>return1</p></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", html, idx)
		t.FailNow()
	}

	input = ">quoted line1.\n\n>```\nprintln(a)\n```\n>return quote\n\n"
	el, idx = parseQuote([]byte(input))
	html = el.Html()
	if html != "<blockquote><p>quoted line1.</p><pre><code>println(a)</code></pre><p>return quote</p></blockquote>" ||
		idx != len(input) {
		t.Log("parseQuote, output:", html, idx)
		t.FailNow()
	}
}

func TestOl(t *testing.T) {
	input := `1. ordered 1
2. ordered 2

>`
	el, idx := parseOl([]byte(input), nil)
	html := el.Html()
	if html != "<ol><li>ordered 1</li><li>ordered 2</li></ol>" ||
		idx != len(input)-1 {
		t.Log("parseOl, output:", html, idx, len(input))
		t.Fail()
	}

	input = `1. ord*er*ed 1
2. ordered [link](test)2`
	el, idx = parseOl([]byte(input), nil)
	html = el.Html()
	if html != `<ol><li>ord<em>er</em>ed 1</li><li>ordered <a href="test">link</a>2</li></ol>` ||
		idx != len(input) {
		t.Log("parseOl, output:", html, idx, len(input))
		t.Fail()
	}

	input = `1. ordered 1

	2. ordered 2

>q`
	el, idx = parseOl([]byte(input), nil)
	html = el.Html()
	if html != `<ol><li><p>ordered 1</p><ol><li>ordered 2</li></ol></li></ol>` ||
		idx != len(input)-2 {
		t.Log("parseOl, output:", html, idx, len(input))
		t.Fail()
	}

	input = `1. >ordered 1
	continue quote

2. ordered 2`
	el, idx = parseOl([]byte(input), nil)
	html = el.Html()
	if html != `<ol><li><blockquote><p>ordered 1 continue quote</p></blockquote></li><li>ordered 2</li></ol>` ||
		idx != len(input) {
		t.Log("parseOl, output:", html, idx)
		t.Fail()
	}
}

func TestNestedUl(t *testing.T) {
	input := `* unordered 1
* unordered 2
	* nested
	* nested2`
	el, idx := parseUl([]byte(input), nil)
	html := el.Html()
	if html != "<ul><li>unordered 1</li><li>unordered 2<ul><li>nested</li><li>nested2</li></ul></li></ul>" ||
		idx != len(input) {
		t.Log("nested, output:", html, idx)
		t.Fail()
	}
}

func TestUlInQuote(t *testing.T) {
	input := `> unordered 1
>
>* unordered 2`
	el, idx := parseQuote([]byte(input))
	html := el.Html()
	if html != "<blockquote><p>unordered 1</p><ul><li>unordered 2</li></ul></blockquote>" ||
		idx != len(input) {
		t.Log("parseUl, output:", html, idx)
		t.Fail()
	}

}

func TestEscape(t *testing.T) {
	input := "im \\*escaped"

	out := Parse([]byte(input))
	str := string(out)
	if str != "<p>im *escaped</p>" {
		//html := el.Html()
		t.Error("out:", str)
	}

}

/*func TestOdd(t *testing.T) {
	input := "`1 + 2 * 3`会先计算乘法，可使用括号`(1 + 2) * 3`，则加法会先计算"

	out := Parse([]byte(input))
	str := string(out)
	if str != "<p>im *escaped</p>" {
		//html := el.Html()
		t.Error("out:", str)
	}
}*/
