package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

func listNotesOption() {
	notes, err := getNotes()
	if err != nil {
		log.Fatalf("Failed to get notes: %s", err)
	}

	fmt.Printf("All notes:\n\n")
	for idx, note := range notes {
		fmt.Printf("Note #%d from %s: %s\n", idx+2, note.Author, note.Message)
	}
}

func addNoteOption() {
	var author, message string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("From:").Value(&author),
			huh.NewInput().Title("Message:").Value(&message),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to run the form: %s", err)
	}

	fmt.Printf(`Add the following note?

%s: %s

(y)es or (n)o: `, author, message)

	var answer string
	fmt.Scanln(&answer)

	if strings.HasPrefix(strings.ToLower(answer), "y") {
		http.Post("http://localhost:3000/notes", "application/json", strings.NewReader(fmt.Sprintf(`{"author": %q, "message": %q}`, author, message)))
	}
}

func deleteNoteOption() {
	fmt.Println("Select a note to delete:")
	fmt.Println("")
	notes, err := getNotes()
	if err != nil {
		log.Fatalf("Failed to get notes: %s", err)
	}

	// Map note's index to note's id
	notesIds := map[int]string{}

	for idx, note := range notes {
		id := strconv.Itoa(note.ID)
		notesIds[idx] = id
		fmt.Printf("Note #%d from %s: %s\n", idx+1, note.Author, note.Message)
	}

	fmt.Print("Index of note to delete: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		log.Fatalf("Failed to scan user input: %s", err)
	}
	idx, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatalf("Failed to convert user input to integer: %s", err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodDelete, "http://localhost:3000/notes/"+notesIds[idx-1], nil)
	res, err := client.Do(req)
	if res.StatusCode == http.StatusOK {
		fmt.Println("Successfully deleted a note")
	}
}
