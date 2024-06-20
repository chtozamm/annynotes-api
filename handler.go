package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := getPosts()

	data, err := json.Marshal(&posts)
	if err != nil {
		log.Fatalf("Failed to marshal posts: %s", err)
	}

	w.Write(data)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	post, err := getPost(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to delete post with id %q: %s", id, err)
	}

	w.Write([]byte(post.Author + ": " + post.Message))
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	post := Post{}

	err := json.Unmarshal(body, &post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to unmarshal request body: %s", err)
	}

	createPost(post)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(post.Author + ": " + post.Message))
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := deletePost(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to delete post with id %q: %s", id, err)
	}
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	newPost := Post{}

	err := json.Unmarshal(body, &newPost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to unmarshal request body: %s", err)
	}

	err = updatePost(id, newPost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to update post with id %q: %s", id, err)
	}
}
