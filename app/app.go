package main

import (
    "fmt"
    "os"

    "lazydap/app/start"
    terminal "lazydap/terminal"

    tea "github.com/charmbracelet/bubbletea"
    // "github.com/charmbracelet/lipgloss"
)

type model struct {
    terminal *terminal.TerminalModel
    start start.StartMenuModel
    width int
    height int
}

func (m model) Init() tea.Cmd {
    var cmds []tea.Cmd
    cmds = append(cmds, m.terminal.Init())
    cmds = append(cmds, m.start.Init())
    return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case error:
        return m, nil
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
            return m, tea.Quit
        case tea.KeyEnter:
            cmds = append(cmds, func() tea.Msg {
                return terminal.RunCmd{
                    Name: "pipenv",
                    Args: []string{"run", "debugpy", "--listen", "5678", "__init__.py"},
                }
            })
        }
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
    }

    m.terminal, cmd = m.terminal.Update(msg)
    cmds = append(cmds, cmd)

    m.start, cmd = m.start.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m model) View() string {
    // message := lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Width(m.width).Height(int(float64(m.height) * 0.7)).Render(
    //     lipgloss.JoinVertical(
    //         lipgloss.Left,
    //         "Input the file name to debug",
    //         m.textInput.View(),
    //         "(ctrl+c to quit)",
    //     ),
    // )
    // terminalMessage := lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Bottom).
    //     Border(lipgloss.RoundedBorder(), true, true).
    //     Render(
    //         m.terminal.View(),
    //     )
    return m.start.View()
    // return lipgloss.Place(
    //     m.width, m.height,
    //     lipgloss.Center,
    //     lipgloss.Center,
    //     lipgloss.JoinVertical(
    //         lipgloss.Center,
    //         message,
    //         terminalMessage,
    //     ),
    // )
}

func initialModel() model {
    return model{
        terminal: terminal.InitialModel(0.99, 0.3),
        start: start.InitialModel(),
    }
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
