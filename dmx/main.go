package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rltvty/go-home/logwrapper"
	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

// The type of our middleware consists of the original handler we want to wrap and a message
type Middleware struct {
	next http.Handler
}

// Make a constructor for our middleware type since its fields are not exported (in lowercase)
func NewMiddleware(next http.Handler) *Middleware {
	return &Middleware{next: next}
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    // WriteHeader(int) is not called if our response implicitly returns 200 OK, so
    // we default to that status code.
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

// Our middleware handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We can modify the request here; for simplicity, we will just log a message
	l := logwrapper.NewLogger()
	defer l.Sync()
	l.APIRequest(r)
	lrw := newLoggingResponseWriter(w)
	m.next.ServeHTTP(lrw, r)
	// We can modify the response here
	l.APIResponse(r, lrw.statusCode)
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)

	log.Fatal(http.ListenAndServe(":8080", NewMiddleware(router)))
}
