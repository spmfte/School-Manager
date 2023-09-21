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
	Editor
)

type model struct {
	currentTab         Tab
	assignments        []string
	selectedAssignment int
	readingMaterials   []ReadingMaterial
	selectedMaterial   int
	notes              []string
	selectedNote       int
	editorContent      string
	originalContent    string
	searchQuery        string
	filteredAssignments []string
}

type ReadingMaterial struct {
	Title  string
	Author string
	Read   bool
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
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
		editorContent:    "Initial file content...",
		originalContent: "Initial file content...",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentTab {
		case Assignments:
			// Handle keyboard navigation for assignments
			switch msg.String() {
			case "ctrl+f":
				m.searchQuery += msg.String()
				m.filteredAssignments = filterAssignments(m.assignments, m.searchQuery)
			default:
				m.searchQuery = ""
				m.filteredAssignments = m.assignments
			}
		case ReadingMaterials:
			// Handle keyboard navigation for reading materials
		case Notes:
			// Handle keyboard navigation for notes
		case Timers:
			// Handle keyboard navigation for timers
		case Editor:
			// Handle text input for the editor
			m.editorContent += msg.String()
		}
		switch msg.String() {
		case "ctrl+left":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "ctrl+right":
			if m.currentTab < Timers {
				m.currentTab++
			}
		case "ctrl+q":
			return m, tea.Quit
		}
	case time.Time:
		// Handle other interactions later
	}
	return m, nil
}

func filterAssignments(assignments []string, query string) []string {
	var results []string
	for _, assignment := range assignments {
		if strings.Contains(assignment, query) {
			results = append(results, assignment)
		}
	}
	return results
}

func (m model) View() string {
	var content string
	switch m.currentTab {
	case Assignments:
		if m.searchQuery != "" {
			content += fmt.Sprintf("Search results for '%s':\n", m.searchQuery)
			content += strings.Join(m.filteredAssignments, "\n")
		} else {
			content += strings.Join(m.assignments, "\n")
		}
	case ReadingMaterials:
		// Render reading materials
	case Notes:
		// Render notes
	case Timers:
		// Render timers
	case Editor:
		content = m.editorContent + "\n\nOriginal Content:\n" + m.originalContent
	}

	// Render the common UI parts (tabs, headers, footers)

	return content
}

