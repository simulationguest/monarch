package main

import (
	"net/http"
)

func main() {
	config := getConfig()
	selector := NewSelector(config)
	http.ListenAndServe(":8081", &selector)
}
