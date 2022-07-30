// Package spin provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/tree/master/spinner
//
// It is useful for displaying that some task is running in the background
// while consuming it's output so that it is not shown to the user.
//
// For example, let's do a long running task: $ sleep 5
//
// We can simply prepend a spinner to this task to show it to the user, while
// performing the task / command in the background.
//
//   $ gum spin -t "Taking a nap..." -- sleep 5
//
// The spinner will automatically exit when the task is complete.
//
package spin

import (
	"os/exec"
    "fmt"
    "os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner spinner.Model
	title   string
	command []string
    display bool
}

type finishCommandMsg struct{ output string }

func commandStart(command []string, display bool) tea.Cmd {
	return func() tea.Msg {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}
		cmd := exec.Command(command[0], args...)
        if display {
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            err := cmd.Run()
            if err != nil{
                fmt.Println(err.Error())
            }
            return finishCommandMsg{output: ""}
        } 
        out, _ := cmd.Output()
		return finishCommandMsg{output: string(out)}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commandStart(m.command, m.display),
	)
}

func (m model) View() string { return m.spinner.View() + " " + m.title }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case finishCommandMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
