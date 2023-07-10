package taskmanager

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"
)

type Task struct {
	command string
	done    *chan Task
	logger  *chan string
	tasker  *chan Task
	cmd     *exec.Cmd
	color   string
	PID     int
}

// func (t *Task) Start() {
// 	if len(t.command) < 2 {
// 		*t.done <- *t
// 		return
// 	}

// 	base := []string{"bash", "-c"}

// 	if runtime.GOOS == "windows" {
// 		base = []string{"cmd", "/C"}
// 	}

// 	cmd := exec.Command(base[0], base[1], t.command)
// 	t.cmd = cmd

// 	regex := regexp.MustCompile(`\033\[[0-9;]*[mK]`)

// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		fmt.Println("Failed to create stderr pipe:", err)
// 		*t.done <- *t
// 		return
// 	}
// 	defer stderr.Close()

// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		fmt.Println("Failed to create stdout pipe:", err)
// 		*t.done <- *t
// 		return
// 	}
// 	defer stdout.Close()

// 	if runtime.GOOS != "windows" {
// 		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// 	}

// 	if err := cmd.Start(); err != nil {
// 		fmt.Println("Failed to start task:", err)
// 		*t.done <- *t
// 		return
// 	}
// 	t.PID, err = t.GetPID()
// 	if err != nil {
// 		fmt.Println("Failed to get process ID:", err)
// 		*t.done <- *t
// 		return
// 	}

// 	if runtime.GOOS != "windows" {
// 		syscall.Setpgid(t.PID, t.PID)
// 	}

// 	*t.tasker <- *t

// 	errscanner := bufio.NewScanner(stderr)
// 	go func() {
// 		defer stderr.Close() // Close the stderr pipe when done
// 		for errscanner.Scan() {
// 			*t.logger <- fmt.Sprintf("[ERROR](bg:red,fg:white)-[%d](fg:%s): %s", t.PID, t.color, regex.ReplaceAllString(errscanner.Text(), ""))
// 		}
// 	}()

// 	scanner := bufio.NewScanner(stdout)
// 	go func() {
// 		defer stdout.Close() // Close the stdout pipe when done
// 		for scanner.Scan() {
// 			*t.logger <- fmt.Sprintf("[%d](fg:%s): %s", t.PID, t.color, regex.ReplaceAllString(scanner.Text(), ""))
// 		}
// 	}()

// 	go func() {
// 		if err := cmd.Wait(); err != nil {
// 			*t.logger <- fmt.Sprintf("[ERROR](bg:red,fg:white)-[%d](fg:%s): %s", t.PID, t.color, err)
// 		}
// 		*t.done <- *t
// 	}()

// 	waitForTerminationSignal()
// 	t.Stop()
// }

// func (task *Task) Stop() error {
// 	err := syscall.Kill(-task.PID, syscall.SIGTERM)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func waitForTerminationSignal() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 	<-c
// }

func (t *Task) Start() {
	if len(t.command) < 2 {
		*t.done <- *t
	}
	base := []string{"bash", "-c"}

	if runtime.GOOS == "windows" {
		base = []string{"cmd", "/C"}
	}

	cmd := exec.Command(base[0], base[1], t.command)
	t.cmd = cmd

	regex := regexp.MustCompile(`\033\[[0-9;]*[mK]`)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	defer stderr.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		*t.done <- *t
	}
	defer stdout.Close()

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if err := cmd.Start(); err != nil {
		*t.done <- *t
	}
	t.PID, err = t.GetPID()
	if err != nil {
		*t.done <- *t
	}
	syscall.Setpgid(t.PID, t.PID)

	*t.tasker <- *t

	errscanner := bufio.NewScanner(stderr)
	go func() {
		defer stderr.Close() // Close the stderr pipe when done
		for errscanner.Scan() {
			*t.logger <- fmt.Sprintf("[ERROR](bg:red,fg:white)-[%d](fg:%s): %s", t.PID, t.color, regex.ReplaceAllString(errscanner.Text(), ""))
		}

	}()
	scanner := bufio.NewScanner(stdout)

	go func() {
		defer stdout.Close() // Close the stdout pipe when done
		for scanner.Scan() {

			*t.logger <- fmt.Sprintf("[%d](fg:%s): %s", t.PID, t.color, regex.ReplaceAllString(scanner.Text(), ""))
		}
	}()

	if err := cmd.Wait(); err != nil {
		*t.logger <- fmt.Sprintf("[ERROR](bg:red,fg:white)-[%d](fg:%s): %s", t.PID, t.color, err)
	}
	*t.done <- *t
}

func (task *Task) Stop() error {
	err := syscall.Kill(-task.PID, syscall.SIGTERM)
	if err != nil {
		return err
	}
	return nil
}

func (t Task) GetPID() (int, error) {
	if t.cmd.Process != nil && t.cmd.Process.Pid > 0 {
		return t.cmd.Process.Pid, nil
	}
	return 0, fmt.Errorf("unable to retrieve command PID")
}
func (t Task) GetColor() string {
	return t.color
}
func (t Task) GetCommands() string {
	return t.command
}
func (t Task) String() string {
	return t.command
}
