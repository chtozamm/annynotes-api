package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

var (
	option string
	// restartForm bool
)

func optionsForm() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What do you want to do?").
				Options(
					huh.NewOption[string]("List all the notes", "1"),
					huh.NewOption[string]("Add new note", "2"),
					huh.NewOption[string]("Update existing note", "3"),
					huh.NewOption[string]("Delete note", "4"),
				).Value(&option),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to run the form: %s", err)
	}

	switch option {
	case "1":
		listNotesOption()
	case "2":
		addNoteOption()
	case "4":
		deleteNoteOption()
	default:
		fmt.Println("Invalid option, exiting.")
		return
	}

	// huh.NewForm(huh.NewGroup(huh.NewConfirm().Affirmative("Continue?").Value(&restartForm))).Run()

	// if restartForm {
	// 	optionsForm()
	// }
}
