# fipmd
[fipmd](https://fipress.org/project/fipmd) is a markdown parser written in Go. It is fast and easy to use. Especially, it provides a set of hooks to let you plant your logic during parsing.

### Usage

**Parse bytes**  
To parse bytes, just use the method `Parse(input []byte) []byte`:

```
out := fipmd.Parse(in)
```

**Parse a string**  
You may use `ParseString(input string) string`.
```
out := fipmd.ParseString(in)
``` 

**Parse a file**
`fipmd` also provides a function to parse a file, `ParseFile(path string) (content []byte, err error)`.
```
content, err := fipmd.ParseFile(path)
```

###Hooks
You may plant your logic to the parser through hooks. The steps is:

1. Create a parser
2. Add hooks to the parser
3. Do parse with the parser

For example:
```
parser := fipmd.NewParser()

parser.AddHeaderHook(func(header *fipmd.Header) {
		println("parsed a header ",header.Text())
	})
parser.AddParaHook(func(para *fipmd.Para) {
		println("parsed a paragraph ",para.Text())
	})

out := parser.Parse(in)
```


