package main

import (
	_ "expvar"
	"net/http"
)

func main() {
	http.ListenAndServe(":1234", nil)
}
