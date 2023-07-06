package main

import (
	"sync"
	"testing"

	taskm "github.com/OmarFaruk-0x01/Go-Commander/taskmanager"
)

func TestNewTask(t *testing.T) {
	wg := sync.WaitGroup{}
	mut := sync.RWMutex{}
	logger := make(chan string)
	tasker := make(chan taskm.Task)
	taskremover := make(chan taskm.Task)
	tm := taskm.NewTaskManager(&wg, &tasker, &taskremover, &logger, &mut)
	go tm.NewTask("echo", "Hello")
	go tm.NewTask("echo", "World")
}
