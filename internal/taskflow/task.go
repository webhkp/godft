package taskflow

import (
	"github.com/webhkp/godft/internal/driver"
)

type TaskI interface {
	GetNextTasks() *[]Task
	AddNextTasks(d *Task)
}

type Task struct {
	Namespace string
	Key       string
	Driver    driver.Driver
	nextTasks []*Task
}

func NewTask(namespace string, key string, driver driver.Driver) (task *Task) {
	task = &Task{
		Namespace: namespace,
		Key:       key,
		Driver:    driver,
	}

	return
}

func (t *Task) GetNextTasks() []*Task {
	return t.nextTasks
}

func (t *Task) AddNextTasks(nextTask *Task) {
	t.nextTasks = append(t.nextTasks, nextTask)
}
