package main

import (
	"fmt"
	"gu"
	"log"
	"net/http"
)

func main() {
	e := gu.New()
	e.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	e.POST("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	log.Fatal(http.ListenAndServe(":9999", e))
}
