package main

import (
	"sync"
	"testing"

	taskm "github.com/OmarFaruk-0x01/Go-Commander/taskmanager"
)

func TestNewTask(t *testing.T) {
	mut := sync.Mutex{}
	logger := make(chan string)
	tasker := make(chan taskm.Task)
	taskremover := make(chan taskm.Task)
	tm := taskm.NewTaskManager(&tasker, &taskremover, &logger, &mut)
	go tm.NewTask("echo", "Hello")
	go tm.NewTask("echo", "World")
}
