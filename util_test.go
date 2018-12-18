package fipmd

import "testing"

func TestExtractTextUntil(t *testing.T) {
	str := "```abcd "

	found, ret, skip := extractTextUntilSpace([]byte(str[4:]), nil)
	t.Log(found, string(ret), skip)
}
