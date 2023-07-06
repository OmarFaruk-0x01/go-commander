package taskmanager

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskManager struct {
	logger      *chan string
	tasker      *chan Task
	taskremover *chan Task
	tasks       []*Task
	wg          *sync.WaitGroup
	mut         *sync.RWMutex
	colors      []string
}

func NewTaskManager(wg *sync.WaitGroup, tasker *chan Task, taskremover *chan Task, logger *chan string, mut *sync.RWMutex) *TaskManager {
	return &TaskManager{
		wg:          wg,
		logger:      logger,
		tasker:      tasker,
		taskremover: taskremover,
		tasks:       []*Task{},
		colors:      []string{"red", "blue", "cyan", "yellow", "green", "white"},
	}
}

func (tm *TaskManager) GetTasks() []*Task {
	return tm.tasks
}

func (tm *TaskManager) NewTask(baseCmd string, args ...string) {
	task := &Task{
		command: append(append([]string{}, baseCmd), args...),
		wg:      tm.wg,
		done:    tm.taskremover,
		logger:  tm.logger,
		tasker:  tm.tasker,
	}
	tm.mut.Lock()
	tm.tasks = append(tm.tasks, task)
	tm.mut.Unlock()
	for i, t := range tm.tasks {
		listLength := len(tm.colors)
		selectedIndex := i % listLength
		if selectedIndex < 0 {
			selectedIndex += listLength - 1
		}
		t.color = tm.colors[selectedIndex]

	}

	task.Start()
}

func (tm *TaskManager) RemoveTask(task *Task) int {
	for i, t := range tm.tasks {
		if t.PID == task.PID {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return i
		}
	}
	return -1
}

func (tm *TaskManager) StopTask(task *Task) {
	task.Stop()
	tm.RemoveTask(task)
	go func(t *Task) {
		*tm.logger <- fmt.Sprintf("[%d](fg:%s): [TERMINATED](bg:red,fg:white)", task.PID, task.GetColor())
	}(task)
}

func (tm *TaskManager) RestartTask(task *Task) {
	task.Stop()
	tm.RemoveTask(task)
	tm.wg.Add(3)
	go func() {
		*tm.logger <- fmt.Sprintf("[%d](fg:%s): [TERMINATED](bg:red,fg:white)", task.PID, task.GetColor())
		tm.wg.Done()
	}()
	go tm.NewTask(task.command[0], task.command...)
	go func() {
		*tm.logger <- fmt.Sprintf("[[%s]](fg:%s): [RESTARTED](bg:cyan,fg:white)", task.String(), task.GetColor())
		tm.wg.Done()
	}()
}

func (tm *TaskManager) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, "Tasks", "Outputs ")
}

func (tm *TaskManager) Init() tea.Cmd {
	// for _, task := range tm.tasks {
	// 	go task.Start()
	// }
	return tea.EnterAltScreen
}

func (tm *TaskManager) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := message.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return tm, tea.Quit
			}
		}
	}

	return tm, nil
}

// func (tm *TaskManager) InitKeyboard() {
// 	err := keyboard.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer keyboard.Close()

// 	for {
// 		ch, _, err := keyboard.GetSingleKey()
// 		if err != nil {
// 			tm.wg.Done()
// 			panic(err)
// 		}
// 		*tm.logger <- strconv.QuoteRuneToASCII(ch)
// 		time.Sleep(time.Millisecond * 100)
// 	}
// }
