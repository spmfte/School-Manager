package main

import (
	"fmt"
	"os"
	"time"
  "strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	Assignments Tab = iota
	ReadingMaterials
	Notes
	Timers
)

type model struct {
	currentTab Tab
	// Introducing state for assignments
	assignments []string
	selectedAssignment int
}

var tabNames = []string{"Assignments", "Reading Materials", "Notes", "Timers"}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v", err)
		os.Exit(1)
	}
}

func initialModel() tea.Model {
	return model{
		assignments: []string{
			"Essay on Shakespeare",
			"Book report on '1984'",
			"Research on Romantic Era",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right":
			if m.currentTab < Timers {
				m.currentTab++
			}
		case "up":
			if m.selectedAssignment > 0 {
				m.selectedAssignment--
			}
		case "down":
			if m.selectedAssignment < len(m.assignments)-1 {
				m.selectedAssignment++
			}
		case "q":
			return m, tea.Quit
		}
	case time.Time:
		// Handle other interactions later
	}
	return m, nil
}

func (m model) View() string {
	// Styles
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDD55")).Height(3).Align(lipgloss.Center)
	tabStyle := lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("#888888"))
	activeTabStyle := tabStyle.Copy().Foreground(lipgloss.Color("#FFDD55")).Underline(true)
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#FFDD55")).Foreground(lipgloss.Color("#000000"))

	// Render tabs
	var tabs []string
	for i, name := range tabNames {
		style := tabStyle
		if m.currentTab == Tab(i) {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(name))
	}

	// Dynamic content rendering based on active tab
	var content string
	switch m.currentTab {
	case Assignments:
		for i, assignment := range m.assignments {
			if i == m.selectedAssignment {
				content += selectedStyle.Render(assignment) + "\n"
			} else {
				content += assignment + "\n"
			}
		}
	case ReadingMaterials:
		content = "List of reading materials will be displayed here."
	case Notes:
		content = "Your notes will be displayed here."
	case Timers:
		content = "Your timers will be displayed here."
	}

	// Combine the layout
	layout := headerStyle.Render("English Class Manager") + "\n\n" +
		strings.Join(tabs, "  ") + "\n\n" +
		content + "\n\n" +
		"Footer: Use arrow keys to navigate tabs. Press 'q' to quit."

	return layout
}

