package main

import "net/http"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := getPosts()
	firstPost := posts[0]

	w.Write([]byte(firstPost.author + ": " + firstPost.message))
	// w.Write([]byte("Bilbo Baggins: Let the adventure begin..."))
}
