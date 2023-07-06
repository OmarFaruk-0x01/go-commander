package taskmanager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Task struct {
	command []string
	wg      *sync.WaitGroup
	done    *chan Task
	logger  *chan string
	tasker  *chan Task
	cmd     *exec.Cmd
	color   string
	PID     int
}

func ErrorLogger(r io.Reader, logger *chan string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		*logger <- scanner.Text()
	}
}

func (t *Task) Start() {
	if len(t.command) < 2 {
		t.wg.Done()
		*t.done <- *t
	}

	cmd := exec.Command(t.command[0], t.command[1:]...)
	t.cmd = cmd

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.wg.Done()
		panic(err)
	}
	defer stderr.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.wg.Done()
		*t.done <- *t
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		t.wg.Done()
		*t.done <- *t
	}
	t.PID, err = t.GetPID()
	if err != nil {
		*t.done <- *t
	}
	*t.tasker <- *t

	t.wg.Add(2)
	go func() {
		defer stderr.Close() // Close the stderr pipe when done
		ErrorLogger(stderr, t.logger)
		t.wg.Done()
	}()
	scanner := bufio.NewScanner(stdout)
	go func() {
		defer stdout.Close() // Close the stdout pipe when done
		for scanner.Scan() {
			*t.logger <- fmt.Sprintf("[%d](fg:%s): %s", t.cmd.Process.Pid, t.color, scanner.Text())
		}
		t.wg.Done()
	}()
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
	t.wg.Done()
	*t.done <- *t
}

func (task *Task) Stop() error {
	process, err := os.FindProcess(task.PID)
	if err != nil {
		return err
	}
	err = process.Signal(os.Interrupt)
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
	str := ""
	for _, c := range t.command[1:] {
		str += " " + c
	}

	return fmt.Sprintf("%s %v", t.cmd.Path, str)
}
