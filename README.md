# Sessions

A minimalist and lightweight HTTP session cookie implementation for Go. Session cookies are encrypted and authenticated using nacl/secretbox.

## Example Usage

```go
package main

import (
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
)

var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

func main() {
	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour

	mux := http.NewServeMux()
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)
	http.ListenAndServe(":4000", session.Enable(mux))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	session.Put(r, "msg", "Hello world")
	w.WriteHeader(200)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	msg := session.GetString(r, "msg")
	w.Write([]byte(msg))
}
```

## TODO

* [docs] Add usage information to the README
* [tests] Test cookie options
* [tests] Improve tests for invalid cookies
* [tests] Increase test coverage
* [tests] Use Travis CI
* [feature] Support multiple named sessions
* [feature] Support Flusher interface
* [feature] Support Hijacker interface