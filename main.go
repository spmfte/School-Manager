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

type Timer struct {
	Description string
	Remaining   time.Duration
}

type model struct {
	currentTab         Tab
	assignments        []string
	selectedAssignment int
	readingMaterials   []ReadingMaterial
	selectedMaterial   int
	notes              []string
	selectedNote       int
	timers             []Timer
	selectedTimer      int
	showModal          bool
	modalContent       string
	inputValue         string
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
		timers: []Timer{
			{Description: "Read Macbeth", Remaining: time.Minute * 30},
			{Description: "Write Essay", Remaining: time.Hour * 2},
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
			switch msg.String() {
			case "esc":
				m.showModal = false
				m.inputValue = ""
			case "enter":
				switch m.currentTab {
				case Assignments:
					m.assignments = append(m.assignments, m.inputValue)
				case ReadingMaterials:
					parts := strings.SplitN(m.inputValue, "-", 2)
					if len(parts) == 2 {
						m.readingMaterials = append(m.readingMaterials, ReadingMaterial{Title: strings.TrimSpace(parts[0]), Author: strings.TrimSpace(parts[1]), Read: false})
					}
				case Notes:
					m.notes = append(m.notes, m.inputValue)
				case Timers:
					parts := strings.SplitN(m.inputValue, "-", 2)
					if len(parts) == 2 {
						duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
						if err == nil {
							m.timers = append(m.timers, Timer{Description: strings.TrimSpace(parts[0]), Remaining: duration})
						}
					}
				}
				m.showModal = false
				m.inputValue = ""
			default:
				m.inputValue += msg.String()
			}
			return m, nil
		}

		// Handle based on active tab
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
				m.modalContent = "Add new assignment:"
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
				m.modalContent = "Add new reading material (Format: Title - Author):"
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
				m.modalContent = "Add new note:"
			}
		case Timers:
			switch msg.String() {
			case "up":
				if m.selectedTimer > 0 {
					m.selectedTimer--
				}
			case "down":
				if m.selectedTimer < len(m.timers)-1 {
					m.selectedTimer++
				}
			case "a":
				m.showModal = true
				m.modalContent = "Set new timer (Format: Description - Duration):"
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
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			for i, name := range tabNames {
				if name == tabNames[m.currentTab] {
					m.currentTab = Tab(i)
					return m, nil
				}
			}
		}
	case time.Time:
		for i := range m.timers {
			if m.timers[i].Remaining > 0 {
				m.timers[i].Remaining -= time.Second
			}
		}
	}
	return m, nil
}

func (m model) View() string {
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

	// Dynamic content rendering
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
		for i, timer := range m.timers {
			line := fmt.Sprintf("%s - %s remaining", timer.Description, timer.Remaining)
			if i == m.selectedTimer {
				content += selectedStyle.Render(line) + "\n"
			} else {
				content += line + "\n"
			}
		}
	}

	if m.showModal {
		modalContent := m.modalContent + "\n" + m.inputValue + "|"
		modal := modalStyle.Render(modalContent)
		content += "\n\n" + modal
	}

	// Combine the layout
	layout := headerStyle.Render("English Class Manager") + "\n\n" +
		strings.Join(tabs, "  ") + "\n\n" +
		content + "\n\n" +
		"Footer: Use arrow keys, mouse wheel, or click to navigate tabs. Press 'q' to quit."

	return layout
}

