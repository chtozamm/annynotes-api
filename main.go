package main

import "net/http"

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", homeHandler)

	http.ListenAndServe("localhost:3000", r)
}
