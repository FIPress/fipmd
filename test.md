## 5.4 Error handling

Go's error handling is accomplished by `error`, `panic` and `recover`.

`error` is used to report an error to the caller, that is the normal error handling case. If there is a fatal, unrecoverable error occur, like `exception` in other languages, `panic` and `recover` can help with it.  

### Errors

`error` a simple built-in interface for representing an error condition. Its zero value is `nil`, which representing no error.

```go
type error interface {
    Error() string
}
```

With Go's multivalue return, it easy to return a detailed error description alongside the normal return value. It is idiomatic to use this style to provide detailed error information. You may find a lot of function defined this way, for example:

```go
func Read(f *File, b []byte) (n int, err error)
```

And a lot of function invocation looks like:

```go
file, err := os.Open("filename")
```

Idiomatically, a library function should return an `error` with necessary information when there is any problem. And the invoker should compare the returned `error` to `nil` first. A `nil` value indicates that no error has occurred and a non-nil value indicates the presence of an error.

Let's have an example. It seems a bit tool long, but actually it's not complicated .

```go
package main

import (
	"encoding/json"
	"errors"
)

var loginFailed = errors.New("login failed")

type argumentError struct {
	argName string
	value   string
	err     string
}

func (ae *argumentError) Error() string {
	return "Argument error:" + ae.err + ", name:" + ae.argName + ", value:" + ae.value
}

type loginData struct {
	Email    string
	Password string
}

func main() {
	inputs := []string{``,
		`{"Email":1,"Password":"b"}`,
		`{"Email":"a","Password":"c"}`,
		`{"Email":"a","Password":"b"}`,
	}

	var err error
	for _, input := range inputs {
		println("login data:", input)

		err = login(input)
		if err != nil {
			if err == loginFailed {
				println("login failed.")
			}

			switch err := err.(type) {
			case *argumentError:
				println("argument error. argument name:", err.argName)
			case *json.UnmarshalTypeError:
				println("unmarshal type error, type:", err.Value, ", expected type:", err.Type.String())
			default:
				println("other kind of error:", err.Error())
			}
		} else {
			println("login succeeded!")
		}

		println("")
	}
}

func login(jsonStr string) (err error) {
	if jsonStr == "" {
		err = &argumentError{argName: "jsonStr", value: jsonStr, err: "should not be empty"}
		return
	}

	ld := new(loginData)
	err = json.Unmarshal([]byte(jsonStr), ld)

	if err != nil {
		return
	}

	if ld.Email == "" || ld.Password == "" {
		err = errors.New("email or Password missing")
		return
	}

	//mock up the matching process
	if ld.Email == "a" && ld.Password != "b" {
		err = loginFailed
		return
	}

	return
}
```

Output:

```
login data: 
argument error. argument name: jsonStr

login data: {"Email":1,"Password":"b"}
unmarshal type error, type: number , expected type: string

login data: {"Email":"a","Password":"c"}
login failed.

login data: {"Email":"a","Password":"b"}
login succeeded!
```

Line 8 defines a variable `loginFailed`, and line 10-18 defines a struct `argumentError`, which implements the built-in interface `error`. The `login` function, line 59-84, shows a function reports errors, and line 36-52 in the `main` function demonstrates how a caller may handle received errors.

To report an error, you may directly create an instance use `errors.New`, like above line 73. Also you may define your own error type, make sure it implements the `error` interface, that is, the `Error() string` method. 

To check an error, first, compare it with `nil`, if non-nil, you may get the information according to the type of the error.

>A library writer should try to make the errors informative and a caller should check the errors and handle them accordingly. Do not ignore errors unless you know exactly what you are doing.

To make an error informative, you may try following ways.

* Define your own struct and provide fields to report error context. As above `argumentError` does. It has `argName` and `value` that the caller can get, as line 45 `err.argName` shows.

* Define your own struct and provide some methods which are return relevant information. For example, the `Err` in the `http` package, it provides a `StatusCode` for the caller to get the http response status code.

```go
type Err struct {
	Response *http.Response
}

...

func (e *Err) StatusCode() int {
	return e.Response.StatusCode
}
```

* A library writer could define a serie of error variables for a caller to compare with, like above `loginFailed`. The caller can directly compare returned error with it as line 40 `if err == loginFailed` shows.

The `http` package also has a lot of errors defined this way.

```go
ErrRepositoryNotFound     = errors.New("repository not found")
ErrEmptyRemoteRepository  = errors.New("remote repository is empty")
ErrAuthenticationRequired = errors.New("authentication required")
ErrAuthorizationFailed    = errors.New("authorization failed")
ErrEmptyUploadPackRequest = errors.New("empty git-upload-pack given")
ErrInvalidAuthMethod      = errors.New("invalid auth method")
ErrAlreadyConnected       = errors.New("session already established")
```
