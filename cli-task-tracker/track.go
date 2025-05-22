package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Task Status Constants
const (
	StatusTodo       = "todo"
	StatusDone       = "done"
	StatusInProgress = "in-progress"
)

// Action Constants
const (
	AddAction      = "add"
	UpdateAction   = "update"
	DeleteAction   = "delete"
	ListAction     = "list"
	DoneAction     = "mark-done"
	ProgressAction = "mark-in-progress"
)

// Task Content
type Task struct {
	ID          uint16 `json:"task-id"`     // No negative numbers
	Status      string `json:"task-status"` // todo/done/in-progress
	Description string `json:"task-description"`
	CreatedAt   string `json:"task-createdAt"` // time.Time
	UpdatedAt   string `json:"task-updatedAt"` // time.Time
}

type TaskList struct {
	Tasks []Task `json:"taskList"`
}

// Meta Data?

func main() {

	// var validActions = []string{AddAction, UpdateAction, DeleteAction, ListAction, DoneAction, ProgressAction}
	// var statuses = []string{StatusTodo, StatusDone, StatusInProgress}
	// Get Action fron Input - Reject Invalid
	// for i := range len(os.Args) {
	// 	if os.Args[1] == validActions[i] {
	// 		fmt.Println("Valid Action: ", os.Args[1])
	// 	} else {
	// 		fmt.Printf("Error: Invalid action. Please use one of the following: %v\n", validActions)
	// 		os.Exit(1)
	// 	}
	// }

	newTask := Task{
		ID:          3,
		Status:      StatusTodo,
		Description: "Sample 3",
		CreatedAt:   getCurrentTimestamp(),
		UpdatedAt:   getCurrentTimestamp(),
	}

	taskList := TaskList{
		Tasks: []Task{newTask},
	}

	WriteTaskToFile("taskList.json", &taskList) // override existing file! Want to append?

}

func WriteTaskToFile(filename string, taskList *TaskList) error {
	fmt.Println("Starting Write Process")

	// 1. Convert struct to JSON
	jsonData, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to convert struct to JSON %w", err)
	}

	// Display JSON
	fmt.Printf("Generated JSON: %s\n", string(jsonData))

	// 2. Create/Check if file exists
	fmt.Printf("Creating/Checking file: %s\n", filename)

	var file *os.File

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File does not exist, creating: %s\n", filename)
		file, err = os.Create(filename)
		if err != nil { // error creating file
			return fmt.Errorf("Failed to create file %w", err)
		}
	} else {
		fmt.Printf("File exists, opening: %s\n", filename) // try to open file
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil { // error opening file
			return fmt.Errorf("Failed to open file %w", err)
		}
	}

	// 3. Write JSON to file
	fmt.Println("Writing JSON to file")
	_, err = file.Write(jsonData)
	if err != nil { // error writing to file
		return fmt.Errorf("Failed to write to file %w", err)
	}

	fmt.Println("JSON written to file successfully")
	defer file.Close() // Close the file after writing
	return nil
}

func getCurrentTimestamp() string {
	return time.Now().Format(time.UnixDate)
}

// func AddTask(description string) {
// 	// Add task to JSON file
// 	// Check if file exists, if not create it
// 	// Append new task to the list
// 	// Update createdAt and updatedAt fields
// }
