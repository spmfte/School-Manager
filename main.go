package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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

type ReadingMaterial struct {
	Title  string
	Author string
	Read   bool
}

type model struct {
	currentTab         Tab
	assignments        []string
	selectedAssignment int
	readingMaterials   []ReadingMaterial
	selectedMaterial   int
	notes              []string
	selectedNote       int
	showModal          bool
	modalContent       string
	// Incorporate timer and other states as we progress
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
		readingMaterials: []ReadingMaterial{
			{Title: "Macbeth", Author: "Shakespeare", Read: true},
			{Title: "1984", Author: "George Orwell", Read: false},
		},
		notes: []string{
			"Note about Macbeth's main theme.",
			"Personal thoughts on 1984.",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showModal {
			if msg.String() == "esc" {
				m.showModal = false
			}
			return m, nil
		}
		// Navigating and selections based on active tab
		switch m.currentTab {
		case Assignments:
			switch msg.String() {
			case "up":
				if m.selectedAssignment > 0 {
					m.selectedAssignment--
				}
			case "down":
				if m.selectedAssignment < len(m.assignments)-1 {
					m.selectedAssignment++
				}
			case "a":
				m.showModal = true
				m.modalContent = "Add new assignment"
			}
		case ReadingMaterials:
			switch msg.String() {
			case "up":
				if m.selectedMaterial > 0 {
					m.selectedMaterial--
				}
			case "down":
				if m.selectedMaterial < len(m.readingMaterials)-1 {
					m.selectedMaterial++
				}
			case "a":
				m.showModal = true
				m.modalContent = "Add new reading material"
			}
		case Notes:
			switch msg.String() {
			case "up":
				if m.selectedNote > 0 {
					m.selectedNote--
				}
			case "down":
				if m.selectedNote < len(m.notes)-1 {
					m.selectedNote++
				}
			case "a":
				m.showModal = true
				m.modalContent = "Add new note"
			}
		}
		// Global navigation
		switch msg.String() {
		case "left":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right":
			if m.currentTab < Timers {
				m.currentTab++
			}
		case "q":
			return m, tea.Quit
		}
	case time.Time:
		// Handle other interactions as needed
	}
	return m, nil
}

func (m model) View() string {
	// Styling
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDD55")).Height(3).Align(lipgloss.Center)
	tabStyle := lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("#888888"))
	activeTabStyle := tabStyle.Copy().Foreground(lipgloss.Color("#FFDD55")).Underline(true)
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#FFDD55")).Foreground(lipgloss.Color("#000000"))
	modalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#333333")).Padding(2)

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
		for i, material := range m.readingMaterials {
			status := "Unread"
			if material.Read {
				status = "Read"
			}
			line := fmt.Sprintf("%s by %s (%s)", material.Title, material.Author, status)
			if i == m.selectedMaterial {
				content += selectedStyle.Render(line) + "\n"
			} else {
				content += line + "\n"
			}
		}
	case Notes:
		for i, note := range m.notes {
			if i == m.selectedNote {
				content += selectedStyle.Render(note) + "\n"
			} else {
				content += note + "\n"
			}
		}
	case Timers:
		content = "Your timers will be displayed here."
	}

	if m.showModal {
		modal := modalStyle.Render(m.modalContent)
		content += "\n\n" + modal
	}

	// Combine the layout
	layout := headerStyle.Render("English Class Manager") + "\n\n" +
		strings.Join(tabs, "  ") + "\n\n" +
		content + "\n\n" +
		"Footer: Use arrow keys to navigate tabs. Press 'q' to quit."

	return layout
}

