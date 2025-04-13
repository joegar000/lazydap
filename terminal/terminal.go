package terminal

import (
    "bufio"
    "fmt"
    "io"
    "os/exec"
    "strings"

    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
)

type TerminalModel struct {
    output []string
    lineCh chan string
    viewport viewport.Model
    widthRatio float64
    heightRatio float64
}

type OutputNewLine struct {
    newLine string
}

type RunCmd struct {
    Name string
    Args []string
}

type KeyMap struct {
    Up key.Binding
    Down key.Binding
}

var DefaultKeyMap = KeyMap{
    Up: key.NewBinding(
        key.WithKeys("k", "up"),        // actual keybindings
        key.WithHelp("↑/k", "move up"), // corresponding help text
    ),
    Down: key.NewBinding(
        key.WithKeys("j", "down"),
        key.WithHelp("↓/j", "move down"),
    ),
}

func InitialModel(widthRatio float64, heightRatio float64) *TerminalModel {
    return &TerminalModel{
        output: []string{},
        lineCh: make(chan string, 1000),
        viewport: viewport.New(0, 0),
        widthRatio: widthRatio,
        heightRatio: heightRatio,
    }
}

// Start the command asynchronously
func startCommand(m *TerminalModel, runCmd RunCmd) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command(runCmd.Name, runCmd.Args...)
        cmd.Dir = "../py/"
        cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
        pipeReader, pipeWriter := io.Pipe()
        cmd.Stdout = pipeWriter
        cmd.Stderr = pipeWriter

        if err := cmd.Start(); err != nil {
            return OutputNewLine{newLine: "Error: " + err.Error()}
        }

        // Start a goroutine to read lines one by one
        go func() {
            scanner := bufio.NewScanner(pipeReader)
            for scanner.Scan() {
                line := scanner.Text()
                m.lineCh <- line
            }
            close(m.lineCh) // Close when done
        }()

        // Return an empty command so Bubble Tea can process the lines asynchronously
        return OutputNewLine{newLine: fmt.Sprintf("Running %s %s", runCmd.Name, strings.Join(runCmd.Args, " "))}
    }
}

// Await the next line of output without blocking
func awaitNextLine(m *TerminalModel) tea.Cmd {
    return func() tea.Msg {
        line, ok := <-m.lineCh
        if ok {
            return OutputNewLine{newLine: line}
        }
        return nil
    }
}

func (m TerminalModel) Init() tea.Cmd {
    return m.viewport.Init()
}

func (m *TerminalModel) Update(msg tea.Msg) (*TerminalModel, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.viewport.Width, m.viewport.Height = int(float64(msg.Width) * m.widthRatio), int(float64(msg.Height) * m.heightRatio)
    case RunCmd:
        cmds = append(cmds, startCommand(m, msg))
    case OutputNewLine:
        m.output = append(m.output, msg.newLine)
        atBottom := m.viewport.AtBottom()
        m.viewport.SetContent(strings.Join(m.output, "\n"))
        if atBottom {
            m.viewport.LineDown(1)
        }
        cmds = append(cmds, awaitNextLine(m))
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, DefaultKeyMap.Up):
            m.viewport.LineUp(1)
        case key.Matches(msg, DefaultKeyMap.Down):
            m.viewport.LineDown(1)
        }
    }
    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m *TerminalModel) View() string {
    return m.viewport.View()
}

