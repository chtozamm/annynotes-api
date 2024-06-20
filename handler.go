package main

import "net/http"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bilbo Baggins: Let the adventure begin..."))
}
