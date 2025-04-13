package start

import (
    config "lazydap/config"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/lipgloss"
)

type StartMenuModel struct {
    width int
    height int
    runConfigs list.Model
}

type ConfigLoaded struct {
    runConfigs []interface{}
}

func InitialModel() StartMenuModel {
    return StartMenuModel{
        width: 0,
        height: 0,
        runConfigs: list.New(),
    }
}

const (
    Finished = "start-menu-complete"
    Pending = "start-menu-pending"
)

func (m StartMenuModel) Init() tea.Cmd {
    return func () tea.Msg {
        config, err := config.EnsureConfigFile()
        if err != nil {
            return err
        }
        return ConfigLoaded{
            runConfigs: config.RunConfigs,
        }
    }
}

func (m StartMenuModel) Update(msg tea.Msg) (StartMenuModel, tea.Cmd) {
    switch msg := msg.(type) {
    case error:
        return m, nil
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
    case ConfigLoaded:
        // TODO: set the items
        m.runConfigs.SetItems()
    }

    return m, nil
}

func createButton(underline string, rest string) string {
    return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render(
        lipgloss.JoinHorizontal(
            lipgloss.Center,
            lipgloss.NewStyle().Underline(true).Render(underline),
            rest,
        ),
    )
}

func (m StartMenuModel) View() string {
    menuWidth := int(float64(m.width) * 0.5)
    menuHeight := int(float64(m.height) * 0.5)
    container := lipgloss.NewStyle().Width(menuWidth).Height(menuHeight).Border(lipgloss.RoundedBorder())

    newButton := createButton("N", "ew")
    runButton := createButton("R", "un")
    editButton := createButton("E", "dit")
    deleteButton := createButton("D", "elete")
    quitButton := createButton("Q", "uit")
    content := lipgloss.NewStyle().Width(menuWidth).Height(menuHeight - lipgloss.Height(newButton)).Render("Getting started")

    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        container.Render(
            lipgloss.JoinVertical(
                lipgloss.Left,
                content,
                lipgloss.JoinHorizontal(
                    lipgloss.Top,
                    newButton,
                    runButton,
                    editButton,
                    deleteButton,
                    quitButton,
                ),
            ),
        ),
    )
}
