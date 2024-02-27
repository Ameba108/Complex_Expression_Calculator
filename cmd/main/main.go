package main

import (
	"complex_expression_calculator/http/server"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("http/templates"))
	http.Handle("/http/templates/", http.StripPrefix("/http/templates/", fs))
	server.HandleRequest()
}
