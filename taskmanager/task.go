package taskmanager

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

type Task struct {
	command []string
	done    *chan Task
	logger  *chan string
	tasker  *chan Task
	cmd     *exec.Cmd
	color   string
	PID     int
}

func (t *Task) Start() {
	if len(t.command) < 2 {
		*t.done <- *t
	}

	cmd := exec.Command(t.command[0], t.command[1:]...)
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

	if err := cmd.Start(); err != nil {
		*t.done <- *t
	}
	t.PID, err = t.GetPID()
	if err != nil {
		*t.done <- *t
	}
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
	process, err := os.FindProcess(task.PID)
	if err != nil {
		return err
	}
	err = process.Signal(syscall.SIGTERM)
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
func (t Task) GetCommands() []string {
	return t.command
}
func (t Task) String() string {
	return strings.Join(t.command, " ")
}
