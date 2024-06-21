package main

type Note struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func main() {
	optionsForm()
}
