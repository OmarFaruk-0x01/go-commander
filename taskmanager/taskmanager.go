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
	mut         *sync.Mutex
	colors      []string
}

func NewTaskManager(tasker *chan Task, taskremover *chan Task, logger *chan string, mut *sync.Mutex) *TaskManager {
	return &TaskManager{
		logger:      logger,
		tasker:      tasker,
		taskremover: taskremover,
		tasks:       []*Task{},
		colors:      []string{"red", "blue", "cyan", "yellow", "green", "black"},
		mut:         mut,
	}
}

func (tm *TaskManager) GetTasks() []*Task {
	return tm.tasks
}

func (tm *TaskManager) NewTask(baseCmd string, args ...string) {
	task := &Task{
		command: append(append([]string{}, baseCmd), args...),
		done:    tm.taskremover,
		logger:  tm.logger,
		tasker:  tm.tasker,
	}
	tm.mut.Lock()
	tm.tasks = append(tm.tasks, task)
	selectedIndex := len(tm.tasks) - 1
	if selectedIndex >= len(tm.colors) {
		selectedIndex = selectedIndex % len(tm.colors)
	}
	task.color = tm.colors[selectedIndex]
	tm.mut.Unlock()

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
	go func() {
		*tm.logger <- fmt.Sprintf("[%d](fg:%s): [TERMINATED](bg:red,fg:white)", task.PID, task.GetColor())
	}()
	go tm.NewTask(task.command[0], task.command[1:]...)
	go func() {
		*tm.logger <- fmt.Sprintf("[[%s]](fg:%s): [RESTARTED](bg:cyan,fg:white)", task.String(), task.GetColor())
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
// 			panic(err)
// 		}
// 		*tm.logger <- strconv.QuoteRuneToASCII(ch)
// 		time.Sleep(time.Millisecond * 100)
// 	}
// }
