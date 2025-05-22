package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices" // Use slices.Contains method
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

// Valid Actions Array
var validActions = []string{
	AddAction,
	UpdateAction,
	DeleteAction,
	ListAction,
	DoneAction,
	ProgressAction,
}

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

// JSON
const fileName = "taskList.json"

func main() {
	// Ensure action provided in arugment
	if len(os.Args) < 2 {
		fmt.Println("Error: No action provided.")
		printUse()
		return
	}

	// Check given action is valid
	action := os.Args[1]
	if !checkAction(action) {
		fmt.Printf("Error: Invalid action. Please use one of the following: %v\n", validActions)
		printUse()
		return
	}

	switch action {
	case AddAction:
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s <description>\n", AddAction)
			return
		}
		add(os.Args[2])

	case UpdateAction:
		if len(os.Args) < 4 {
			fmt.Printf("Usage: %s <task-id> <description>\n", UpdateAction)
			return
		}
		update(os.Args[2], os.Args[3])

	case DeleteAction:
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s <task-id>\n", DeleteAction)
			return
		}
		delete(os.Args[2])

	case ListAction:
		list(os.Args)

	case DoneAction:
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s <task-id>\n", DoneAction)
			return
		}
		markDone(os.Args[2])

	case ProgressAction:
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s <task-id>\n", ProgressAction)
			return
		}
		markInProgress(os.Args[2])
	}

}

// Load Data
func load() TaskList {
	var tList TaskList
	data, err := os.ReadFile(fileName)
	if err == nil {
		json.Unmarshal(data, &tList)
	}
	return tList
}

// Save Data
func save(tList TaskList) {
	data, _ := json.MarshalIndent(tList, "", "  ")
	os.WriteFile(fileName, data, 0644)
}

// Action Methods
func add(description string) {
	tList := load()
	var maxID uint16
	for _, t := range tList.Tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	tList.Tasks = append(tList.Tasks, Task{
		ID:          maxID + 1,
		Status:      "todo",
		Description: description,
		CreatedAt:   getCurrentTimestamp(),
		UpdatedAt:   getCurrentTimestamp(),
	})
	save(tList)
	fmt.Println("Task Added: ", tList.Tasks[len(tList.Tasks)-1].ID) // Print the ID of the added task
}

func update(idStr, description string) {
	tList := load()
	for i, t := range tList.Tasks {
		if fmt.Sprintf("%d", t.ID) == idStr {
			tList.Tasks[i].Description = description
			tList.Tasks[i].UpdatedAt = getCurrentTimestamp()
			save(tList)
			fmt.Println("Task updated")
			return
		}
	}
	fmt.Println("Task not found")
}

func delete(idStr string) {
	tList := load()
	for i, t := range tList.Tasks {
		if fmt.Sprintf("%d", t.ID) == idStr {
			tList.Tasks = slices.Delete(tList.Tasks, i, i+1)
			save(tList)
			fmt.Println("Task deleted")
			return
		}
	}
	fmt.Println("Task not found")
}

// List Tasks - Filter by Status on Args
func list(args []string) {
	tList := load()
	if len(tList.Tasks) == 0 {
		fmt.Println("No tasks")
		return
	}
	if len(args) < 3 { // No extra args - List all Tasks
		for _, t := range tList.Tasks {
			fmt.Printf("[%d] %s - %s\n", t.ID, t.Status, t.Description)
		}
		return
	}

	var status string // Get status to filter
	switch args[2] {
	case StatusDone:
		status = StatusDone
	case StatusInProgress:
		status = StatusInProgress
	case StatusTodo:
		status = StatusTodo
	default:
		fmt.Printf("Unknown status: %s\n", args[2])
		return
	}
	for _, t := range tList.Tasks {
		if t.Status == status {
			fmt.Printf("[%d] %s - %s\n", t.ID, t.Status, t.Description)
		}
	}
}

func markDone(idStr string) {
	tList := load()
	for i, t := range tList.Tasks {
		if fmt.Sprintf("%d", t.ID) == idStr {
			tList.Tasks[i].Status = StatusDone
			tList.Tasks[i].UpdatedAt = getCurrentTimestamp()
			save(tList)
			fmt.Println("Task marked as done")
			return
		}
	}
	fmt.Println("Task not found")
}

func markInProgress(idStr string) {
	tList := load()
	for i, t := range tList.Tasks {
		if fmt.Sprintf("%d", t.ID) == idStr {
			tList.Tasks[i].Status = StatusInProgress
			tList.Tasks[i].UpdatedAt = getCurrentTimestamp()
			save(tList)
			fmt.Println("Task marked as in-progress")
			return
		}
	}
	fmt.Println("Task not found")
}

// Helper and Print Functions
func checkAction(action string) bool {
	return slices.Contains(validActions, action)
}

// Get Nice Timestamp
func getCurrentTimestamp() string {
	return time.Now().Format(time.UnixDate)
}

// Display usage with Actions
func printUse() {
	fmt.Println("Task Tracker CLI")
	fmt.Println("Usage: cli-task-tracker <action> [options]")
	for _, action := range validActions {
		switch action {
		case AddAction:
			fmt.Println("  add <task-description>")
		case UpdateAction:
			fmt.Println("  update <task-id> <task-description>")
		case DeleteAction:
			fmt.Println("  delete <task-id>")
		case ListAction:
			fmt.Println("  list")
		case DoneAction:
			fmt.Println("  mark-done <task-id>")
		case ProgressAction:
			fmt.Println("  mark-in-progress <task-id>")
		}
	}
}
