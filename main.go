// // // package main

// // // import (
// // // 	"bufio"
// // // 	"bytes"
// // // 	"fmt"
// // // 	"io"
// // // 	"log"
// // // 	"os"
// // // 	"os/exec"
// // // 	"sync"
// // // 	"time"

// // // 	"github.com/eiannone/keyboard"
// // // )

// // // type Task struct {
// // // 	Command string
// // // 	PID     int
// // // 	cmd     *exec.Cmd
// // // 	Logs    *bytes.Buffer
// // // 	done    chan struct{}
// // // 	wg      sync.WaitGroup
// // // }

// // // type TaskManager struct {
// // // 	tasks     []*Task
// // // 	taskCount int
// // // }

// // // func (task *Task) Start() {
// // // 	// Create a pipe for capturing the task's output
// // // 	stdoutPipe, _ := task.cmd.StdoutPipe()
// // // 	stderrPipe, _ := task.cmd.StderrPipe()

// // // 	// Start the task
// // // 	if err := task.cmd.Start(); err != nil {
// // // 		log.Println("Error: ", err)
// // // 	}

// // // 	// Create a scanner for reading the output
// // // 	scanner := bufio.NewScanner(io.MultiReader(stdoutPipe, stderrPipe))
// // // 	for scanner.Scan() {
// // // 		// Read each line of output and write it to the task's logs buffer
// // // 		task.Logs.WriteString(scanner.Text() + "\n")
// // // 	}

// // // 	// Signal that the task has completed
// // // 	task.done <- struct{}{}
// // // 	task.wg.Done()
// // // }

// // // func (task *Task) Stop() {
// // // 	// Stop the task gracefully
// // // 	process, err := os.FindProcess(task.PID)
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}
// // // 	err = process.Signal(os.Interrupt)
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}

// // // 	// Wait for the task to complete
// // // 	<-task.done
// // // 	task.wg.Wait()
// // // }

// // // func NewTaskManager() *TaskManager {
// // // 	return &TaskManager{}
// // // }

// // // func (tm *TaskManager) StartTask(command string) {
// // // 	task := &Task{
// // // 		Command: command,
// // // 		Logs:    &bytes.Buffer{},
// // // 		done:    make(chan struct{}),
// // // 	}

// // // 	cmd := exec.Command("bash", "-c", command)
// // // 	task.cmd = cmd
// // // 	log.Println("Start Task", cmd)
// // // 	pid, err := getCmdPID(cmd)
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}
// // // 	task.PID = pid

// // // 	tm.tasks = append(tm.tasks, task)
// // // 	tm.taskCount++

// // // 	task.wg.Add(1)
// // // 	go task.Start()
// // // }

// // // func (tm *TaskManager) RestartTask(pid int, command string) {
// // // 	for _, task := range tm.tasks {
// // // 		if task.PID == pid {
// // // 			task.Stop()
// // // 			tm.removeTask(task)
// // // 			tm.StartTask(command)
// // // 			return
// // // 		}
// // // 	}
// // // 	// Task not found, handle error or notify user
// // // }

// // // func (tm *TaskManager) removeTask(task *Task) {
// // // 	for i, t := range tm.tasks {
// // // 		if t == task {
// // // 			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
// // // 			tm.taskCount--
// // // 			return
// // // 		}
// // // 	}
// // // }

// // // func (tm *TaskManager) ShowPrompt() {
// // // 	fmt.Print("Enter the task number to restart or kill (or press 'Q' to quit): ")
// // // 	char, _, err := keyboard.GetSingleKey()
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}

// // // 	if char == 'Q' || char == 'q' {
// // // 		fmt.Println("Exiting...")
// // // 		os.Exit(0)
// // // 	}

// // // 	taskNumber := int(char - '0')
// // // 	if taskNumber < 0 || taskNumber >= tm.taskCount {
// // // 		// Invalid task number, handle error or notify user
// // // 		fmt.Println("Invalid task number!")
// // // 		return
// // // 	}

// // // 	fmt.Println("Selected task number:", taskNumber)

// // // 	fmt.Print("Enter 'R' to restart or 'K' to kill the task: ")
// // // 	char, _, err = keyboard.GetSingleKey()
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}

// // // 	switch char {
// // // 	case 'R', 'r':
// // // 		fmt.Println("Restarting task...")
// // // 		tm.RestartTask(tm.tasks[taskNumber].PID, tm.tasks[taskNumber].Command)
// // // 	case 'K', 'k':
// // // 		fmt.Println("Killing task...")
// // // 		tm.StopTask(tm.tasks[taskNumber].PID)
// // // 	default:
// // // 		// Invalid key, handle error or notify user
// // // 		fmt.Println("Invalid key!")
// // // 	}
// // // }

// // // func (tm *TaskManager) ShowLogs() {
// // // 	go func() {
// // // 		for i, task := range tm.tasks {
// // // 			fmt.Printf("[%d] Logs for task with PID %d:\n%s\n", i, task.PID, task.Logs.String())
// // // 		}
// // // 	}()

// // // }

// // // func getCmdPID(cmd *exec.Cmd) (int, error) {
// // // 	if cmd.Process != nil && cmd.Process.Pid > 0 {
// // // 		return cmd.Process.Pid, nil
// // // 	}
// // // 	return 0, fmt.Errorf("unable to retrieve command PID")
// // // }

// // // func (tm TaskManager) initKeyboard() {
// // // 	err := keyboard.Open()
// // // 	if err != nil {
// // // 		// Handle error
// // // 	}
// // // 	defer keyboard.Close()

// // // 	for {
// // // 		tm.ShowLogs()
// // // 		tm.ShowPrompt()
// // // 		time.Sleep(time.Millisecond * 100)
// // // 	}
// // // }

// // // func Run() {
// // // 	tm := NewTaskManager()

// // // 	// Example usage:
// // // 	tm.StartTask("ping google.com")
// // // 	tm.StartTask("ls -l")
// // // 	tm.StartTask("go version")

// // // 	go tm.initKeyboard()
// // // 	// Wait indefinitely
// // // 	select {}
// // // }

// // // func main() {
// // // 	Run()
// // // }

package main

import (
	"fmt"
	"log"
	"sync"

	taskm "github.com/OmarFaruk-0x01/Go-Commander/taskmanager"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type State struct {
	selectedTask int
	selectedMenu int
	focusedTab   string
	isDialogOpen bool

	tm *taskm.TaskManager
}

var (
	focusedTabStyle   = ui.NewStyle(ui.ColorMagenta)
	normalTabStyle    = ui.NewStyle(ui.ColorWhite)
	selectedItemStyle = ui.NewStyle(ui.ColorCyan)
)

func (s State) ViewDialog() (*ui.Grid, *widgets.List) {
	w, h := ui.TerminalDimensions()
	grid := ui.NewGrid()
	grid.Border = true
	grid.SetRect(0, 0, w, h)
	p2 := widgets.NewParagraph()
	p2.Text = "Down:[<\u21A7>](fg:red)  Up:[<\u21A5>](fg:red) Select:[<\u21A6>](fg:red)"
	p2.SetRect(0, 0, 0, 0)
	menu := widgets.NewList()
	menu.Title = fmt.Sprintf("Setting [%s]", s.tm.GetTasks()[s.selectedTask])
	menu.BorderStyle = ui.NewStyle(ui.ColorBlue)
	menu.Rows = append(menu.Rows,
		"[1] Restart Task",
		"[2] Terminate Task",
	)
	menu.SelectedRowStyle = selectedItemStyle
	grid.Set(ui.NewRow(1.9/2, menu), ui.NewRow(.1/2, p2))
	return grid, menu
}

func main() {
	wg := sync.WaitGroup{}
	mut := sync.RWMutex{}
	logger := make(chan string)
	tasker := make(chan taskm.Task)
	taskremover := make(chan taskm.Task)
	tm := taskm.NewTaskManager(&wg, &tasker, &taskremover, &logger, &mut)
	state := State{
		selectedTask: 0,
		selectedMenu: -1,
		focusedTab:   "task",
		isDialogOpen: false,
		tm:           tm,
	}

	if err := ui.Init(); err != nil {
		log.Fatal("failed to initialized window")
	}
	defer ui.Close()

	wg.Add(4)
	go tm.NewTask("ping", "facebook.com")
	go tm.NewTask("ping", "google.com")
	go tm.NewTask("echo", "hello")
	go tm.NewTask("echo", "world")

	termWidth, termHeight := ui.TerminalDimensions()
	dialog, menu := state.ViewDialog()

	logList := widgets.NewList()
	logList.Title = "Logs"

	help := widgets.NewList()
	help.Title = "Info"

	table := widgets.NewTable()
	table.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
	table.RowSeparator = false
	table.TextAlignment = ui.AlignCenter
	table.FillRow = true
	table.BorderStyle = focusedTabStyle
	table.Rows = [][]string{
		{"NO", "PID", "Command"},
	}
	table.ColumnWidths = []int{5, 15, int(float64(termWidth)*(0.25/2.0)) + 6}

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, table), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
	)

	ui.Render(grid)
	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Enter>":
				if state.focusedTab == "setting" && state.isDialogOpen {
					if menu.SelectedRow == 0 {
						// restart action
						tm.RestartTask(tm.GetTasks()[state.selectedTask-1])
					} else {
						tm.StopTask(tm.GetTasks()[state.selectedTask-1])
					}
					state.isDialogOpen = false
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, table), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "task"
					table.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle
					ui.Render(grid)
				}
			case "<Escape>":
				if state.isDialogOpen {
					state.isDialogOpen = false
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, table), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "task"
					table.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle
					ui.Render(grid)
				}
			case "<Space>":
				if !state.isDialogOpen && state.selectedTask != 0 {
					state.isDialogOpen = true
					grid.Items = make([]*ui.GridItem, 0)
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, table), ui.NewRow(1.0/2, dialog)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "setting"
					table.BorderStyle = normalTabStyle
					menu.BorderStyle = focusedTabStyle
					menu.Title = fmt.Sprintf("Settings [%s]", tm.GetTasks()[state.selectedTask-1].String())
					ui.Render(grid)
				} else {
					state.isDialogOpen = false
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, table), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
					)
					ui.Render(table)
				}
			case "<Tab>":
				switch state.focusedTab {
				case "task":
					state.focusedTab = "log"
					table.BorderStyle = normalTabStyle
					logList.BorderStyle = focusedTabStyle
				case "log":
					state.focusedTab = "task"
					table.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle

				}
				ui.Render(grid)
			case "<Down>":
				switch state.focusedTab {
				case "task":
					state.selectedTask++
					if state.selectedTask > len(table.Rows)-1 {
						state.selectedTask = 1
					}
					for i := 1; i < len(table.Rows); i++ {
						if i == state.selectedTask {
							table.RowStyles[i] = selectedItemStyle
						} else {
							table.RowStyles[i] = normalTabStyle
						}
					}
				case "setting":
					menu.SelectedRow++
					if menu.SelectedRow > len(menu.Rows)-1 {
						menu.SelectedRow = 0
					}
					ui.Render(menu)
				}
				ui.Render(table)
			case "<Up>":
				switch state.focusedTab {
				case "task":
					state.selectedTask--
					if state.selectedTask < 1 {
						state.selectedTask = len(table.Rows) - 1
					}
					for i := 0; i < len(table.Rows); i++ {
						if i == state.selectedTask {
							table.RowStyles[i] = selectedItemStyle
						} else {
							table.RowStyles[i] = normalTabStyle
						}
					}
				case "setting":
					menu.SelectedRow--
					if menu.SelectedRow < 0 {
						menu.SelectedRow = len(menu.Rows) - 1
					}
					ui.Render(menu)
				}
				ui.Render(table)
			case "j", "<MouseWheelDown>":
				logList.ScrollDown()
				ui.Render(logList)
			case "k", "<MouseWheelUp>":
				logList.ScrollUp()
				ui.Render(logList)
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			default:
				log.Printf("UiEvents: %v", e)
			}
		case <-tasker:
			{
				for i, v := range tm.GetTasks() {
					table.Rows = append(table.Rows, []string{fmt.Sprintf("%d", i), fmt.Sprintf("[%d](fg:%s)", v.PID, v.GetColor()), v.String()})

				}
				table.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
				ui.Render(table)
			}
		case doneTask := <-taskremover:
			{
				tm.RemoveTask(&doneTask)
				table.Rows = table.Rows[:1]
				for i, v := range tm.GetTasks() {
					table.Rows = append(table.Rows, []string{fmt.Sprintf("%d", i), fmt.Sprintf("[%d](fg:%s)", v.PID, v.GetColor()), v.String()})
				}
				table.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
				ui.Render(table)
			}
		case log := <-logger:
			{
				logList.Rows = append(logList.Rows, log)
				logList.ScrollDown()
				ui.Render(logList)
			}
		}

	}
}

// package main

// import (
// 	"log"
// 	"sync"

// 	"github.com/OmarFaruk-0x01/Go-Commander/taskmanager"
// 	tea "github.com/charmbracelet/bubbletea"
// )

// func main() {
// 	wg := sync.WaitGroup{}
// 	logger := make(chan string)
// 	tasker := make(chan taskmanager.Task)
// 	taskremover := make(chan taskmanager.Task)
// 	tm := taskmanager.NewTaskManager(&wg, &tasker, &taskremover, &logger)

// 	// Register tasks
// 	go tm.NewTask("echo", "Task added")
// 	go tm.NewTask("ping", "facebook.com")
// 	go tm.NewTask("ping", "google.com")

// 	p := tea.NewProgram(tm, tea.WithAltScreen())

// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
