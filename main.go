package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
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
	autoScroll   bool

	taskPage      int
	taskListLimit int

	tm *taskm.TaskManager
}
type Cmds []string

func (c *Cmds) String() string {
	return strings.Join(*c, " ")
}
func (c *Cmds) Set(value string) error {
	*c = append(*c, value)
	return nil
}

var (
	focusedTabStyle   = ui.NewStyle(ui.ColorMagenta)
	normalTabStyle    = ui.NewStyle(ui.ColorWhite)
	selectedItemStyle = ui.NewStyle(ui.ColorCyan)
)

func (s State) ViewDialog() (*ui.Grid, *widgets.List) {
	grid := ui.NewGrid()
	grid.Border = true
	p2 := widgets.NewParagraph()
	p2.Text = "Down:[<\u21A7>](fg:red)  Up:[<\u21A5>](fg:red) Select:[<\u21A6>](fg:red)"
	p2.SetRect(0, 0, 0, 0)
	menu := widgets.NewList()
	menu.BorderStyle = ui.NewStyle(ui.ColorBlue)
	menu.Rows = append(menu.Rows,
		"[1] Restart Task",
		"[2] Terminate Task",
		"[3] Go back",
	)
	menu.SelectedRowStyle = selectedItemStyle
	grid.Set(ui.NewRow(1.9/2, menu), ui.NewRow(.1/2, p2))
	return grid, menu
}

func main() {
	var mut sync.Mutex
	var cmds Cmds
	logger := make(chan string)
	tasker := make(chan taskm.Task)
	taskremover := make(chan taskm.Task)
	tm := taskm.NewTaskManager(&tasker, &taskremover, &logger, &mut)
	state := State{
		selectedTask:  0,
		selectedMenu:  -1,
		focusedTab:    "task",
		isDialogOpen:  false,
		taskPage:      1,
		taskListLimit: 5,
		tm:            tm,
		autoScroll:    true,
	}

	flag.Var(&cmds, "cmd", "shell commands \"<cmd>\" (e.g. ./commander -cmd \"echo hello\" -cmd \"echo world\" ...)")
	flag.Parse()

	if len(cmds) == 0 {
		flag.PrintDefaults()
		return
	}

	if err := ui.Init(); err != nil {
		log.Fatal("failed to initialized window")
	}
	defer ui.Close()

	for _, c := range cmds {
		cmd := strings.Split(c, " ")
		go tm.NewTask(cmd[0], cmd[1:]...)
	}

	termWidth, termHeight := ui.TerminalDimensions()
	dialog, menu := state.ViewDialog()

	logList := widgets.NewList()
	logList.Title = "Logs"
	logList.SelectedRowStyle = ui.NewStyle(ui.ColorClear)

	help := widgets.NewList()
	help.Title = "Info"
	help.WrapText = true
	help.Rows = []string{
		"Key Bindings",
		"---------------------------------------------------------",
		"[<q>/<ctrl+c>](fg:yellow): 	Kill all Tasks and exit",
		"[<tab>](fg:yellow):			Move focus from tabs",
		"[<up>](fg:yellow): 			Navigate Up in focused tab",
		"[<down>](fg:yellow): 			Navigate Down in focused tab",
		"[<space>](fg:yellow): 			Enter settings of selected item",
		"[<esc>](fg:yellow): 			Go back from settings",
		"[<enter>](fg:yellow): 			Execute Action",
		"[<scroll-up>](fg:yellow): 		Navigate Logs upside & (autoscroll: off)",
		"[<scroll-down>](fg:yellow): 	Navigate Logs upside & (autoscroll: off)",
		"[<ctrl+s>](fg:yellow): 		Toggle AutoScroll Logs",
	}

	taskList := widgets.NewList()
	taskList.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
	// taskList.RowSeparator = false
	// taskList.TextAlignment = ui.AlignCenter
	// taskList.FillRow = true
	taskList.BorderStyle = focusedTabStyle
	taskList.SelectedRowStyle = normalTabStyle
	taskList.WrapText = true
	taskList.Rows = []string{
		"[ PID ] - Task",
	}
	// taskList.ColumnWidths = []int{5, 15, int(float64(termWidth)*(0.25/2.0)) + 6}

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, taskList), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
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
						tm.RestartTask(tm.GetTasks()[taskList.SelectedRow-1])
					} else if menu.SelectedRow == 1 {
						tm.StopTask(tm.GetTasks()[taskList.SelectedRow-1])
					}
					state.isDialogOpen = false
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, taskList), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "task"
					taskList.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle
					ui.Render(grid)
					ui.Render(help)
				}
			case "<Escape>":
				if state.isDialogOpen {
					state.isDialogOpen = false
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, taskList), ui.NewRow(1.0/2, help)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "task"
					taskList.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle
					ui.Render(help)
					ui.Render(taskList)
				}
			case "<Space>":
				if !state.isDialogOpen && taskList.SelectedRow != 0 && state.focusedTab == "task" {
					state.isDialogOpen = true
					grid.Items = make([]*ui.GridItem, 0)
					grid.Set(
						ui.NewRow(1.0, ui.NewCol(.5/2, ui.NewRow(1.0/2, taskList), ui.NewRow(1.0/2, dialog)), ui.NewCol(1.5/2, logList)),
					)
					state.focusedTab = "setting"
					taskList.BorderStyle = normalTabStyle
					menu.BorderStyle = focusedTabStyle
					menu.Title = fmt.Sprintf("Settings [%s]", tm.GetTasks()[taskList.SelectedRow-1].String())
					ui.Render(grid)
				}
			case "<Tab>":
				switch state.focusedTab {
				case "task":
					state.focusedTab = "logs"
					taskList.BorderStyle = normalTabStyle
					logList.BorderStyle = focusedTabStyle
				case "logs":
					state.focusedTab = "task"
					taskList.BorderStyle = focusedTabStyle
					logList.BorderStyle = normalTabStyle

				}
				ui.Render(grid)
			case "<Down>":
				switch state.focusedTab {
				case "task":
					taskList.SelectedRow++
					if taskList.SelectedRow > len(taskList.Rows)-1 {
						taskList.SelectedRow = 0
					}
					if taskList.SelectedRow == 0 {
						taskList.SelectedRowStyle = normalTabStyle
					} else {
						taskList.SelectedRowStyle = selectedItemStyle
					}
				// for i := 1; i < len(taskList.Rows); i++ {
				// 	if i == taskList.SelectedRow {
				// 		taskList.RowStyles[i] = selectedItemStyle
				// 	} else {
				// 		taskList.RowStyles[i] = normalTabStyle
				// 	}
				// }
				case "logs":
					logList.ScrollDown()
					ui.Render(taskList)
				case "setting":
					menu.SelectedRow++
					if menu.SelectedRow > len(menu.Rows)-1 {
						menu.SelectedRow = 0
					}
					ui.Render(menu)
				}
				ui.Render(taskList)
			case "<Up>":
				switch state.focusedTab {
				case "task":
					taskList.SelectedRow--
					if taskList.SelectedRow < 0 {
						taskList.SelectedRow = len(taskList.Rows) - 1
					}
					if taskList.SelectedRow == 0 {
						taskList.SelectedRowStyle = normalTabStyle
					} else {
						taskList.SelectedRowStyle = selectedItemStyle
					}
					// for i := 0; i < len(taskList.Rows); i++ {
					// 	if i == taskList.SelectedRow {
					// 		taskList.RowStyles[i] = selectedItemStyle
					// 	} else {
					// 		taskList.RowStyles[i] = normalTabStyle
					// 	}
					// }
					ui.Render(taskList)

				case "logs":
					logList.ScrollUp()
					ui.Render(taskList)
				case "setting":
					menu.SelectedRow--
					if menu.SelectedRow < 0 {
						menu.SelectedRow = len(menu.Rows) - 1
					}
					ui.Render(menu)
				}
			case "j", "<MouseWheelDown>":
				state.autoScroll = false
				logList.Title = "Logs [Auto-Scroll: Off]"
				// logList.ScrollAmount(10)
				logList.ScrollDown()
				ui.Render(logList)
			case "k", "<MouseWheelUp>":
				state.autoScroll = false
				logList.Title = "Logs [Auto-Scroll: Off]"
				// logList.ScrollAmount(10)
				logList.ScrollUp()
				ui.Render(logList)
			case "q", "<C-c>":
				for _, t := range tm.GetTasks() {
					t.Stop()
				}
				return
			case "<C-s>":
				state.autoScroll = !state.autoScroll
				if state.autoScroll {
					logList.Title = "Logs [Auto-Scroll: On]"
				} else {
					logList.Title = "Logs [Auto-Scroll: Off]"
				}
				logList.ScrollPageDown()
				ui.Render(logList)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-tasker:
			{
				taskList.Rows = taskList.Rows[:1]
				for _, v := range tm.GetTasks() {
					taskList.Rows = append(taskList.Rows, fmt.Sprintf("[%d](fg:%s) - %s", v.PID, v.GetColor(), v.String()))

				}
				taskList.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
				ui.Render(taskList)
			}
		case doneTask := <-taskremover:
			{
				tm.RemoveTask(&doneTask)
				taskList.Rows = taskList.Rows[:1]
				for _, v := range tm.GetTasks() {
					taskList.Rows = append(taskList.Rows, fmt.Sprintf("[%d](fg:%s) - %s", v.PID, v.GetColor(), v.String()))
				}
				taskList.Title = fmt.Sprintf("Tasks [%d]", len(tm.GetTasks()))
				ui.Render(taskList)
			}
		case log := <-logger:
			{
				logList.Rows = append(logList.Rows, log)
				if state.autoScroll {
					logList.ScrollPageDown()
				}
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
