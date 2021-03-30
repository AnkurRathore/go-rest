package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

//Taskstore will be an in memory store that can be accessed concurrently

type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

//CreateTask creates a new Task
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

//GetTask retrieves a task from the store by id and if it does not
// exist an error is returned
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}

}

//DeleteTask  deletes a Task by ID
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

//DeleteAllTasks deletes all the tasks in the store
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)

	return nil
}

//GetAllTasks returns all the tasks in the store
func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

//GetTaskByTag returns the tasks that match the tag
func (ts *TaskStore) GetTaskByTag(tag string) []Task {
	ts.Lock()

	defer ts.Unlock()

	var tasks []Task

taskloop:
	for _, task := range ts.tasks {
		for _, tasktag := range task.Tags {
			if tasktag == tag {
				tasks = append(tasks, task)
				continue taskloop
			}
		}
	}

	return tasks
}

//GetTaskByDueDate returns all the tasks that have the give due date
func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()

	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()

		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
