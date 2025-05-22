# CLI Task Tracker

**From** Roadmap.sh - Task Tracker
https://roadmap.sh/projects/task-tracker

Completed task is `track.go` that loads/checks `taskList.json`

## Task Requirements

Run from command line, accept user actions/inputs and store task in JSON file:

* Add, update, delete tasks
* Mark task as 'in progress' or 'done'
* List all tasks
* List all tasks done
* List all tasks not done
* List all tasks in progress

### Notes

Positional arguments to accept inputs
JSON file to store tasks in current directory
JSON file should be created if does not exist
Native file system module in lang
No external libraries
Handle errors/edge cases propoerly.

### Example

```
# Adding a new task
task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Updating and deleting tasks
task-cli update 1 "Buy groceries and cook dinner"
task-cli delete 1

# Marking a task as in progress or done
task-cli mark-in-progress 1
task-cli mark-done 1

# Listing all tasks
task-cli list

# Listing tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress
```

### Task Properties

Each task should have the following properties:

* `id`: A unique identifier for the task
* `description`: A short description of the task
* `status`: The status of the task (todo, in-progress, done)
* `createdAt`: The date and time when the task was created
* `updatedAt`: The date and time when the task was last updated

Make sure to add these properties to the JSON file when adding a new task and update them when updating a task.
