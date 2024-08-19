package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mebn/termfm/internal/audioplayer"
)

type countryItem struct {
	title, desc string
}

func (i countryItem) Title() string       { return i.title }
func (i countryItem) Description() string { return i.desc }
func (i countryItem) FilterValue() string { return i.title }

type stationItem struct {
	title, desc string
	url         string
}

func (i stationItem) Title() string       { return i.title }
func (i stationItem) Description() string { return i.desc }
func (i stationItem) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("TermFM")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "right", "l":
			if m.state == countryState {
				m.state = stationState
			} else {
				m.state = countryState
			}

		case "left", "h":
			if m.state == countryState {
				m.state = stationState
			} else {
				m.state = countryState
			}

		case "enter":
			if m.state == stationState {
				station := m.stationsList.SelectedItem().(stationItem)
				go audioplayer.PlayStation(station.url)
			}
		}

	case tea.WindowSizeMsg:
		m.countriesList.SetSize(msg.Width, msg.Height)
		m.stationsList.SetSize(msg.Width, msg.Height)
	}

	var listCmd tea.Cmd
	if m.state == countryState {
		m.countriesList, listCmd = m.countriesList.Update(msg)
	} else if m.state == stationState {
		m.stationsList, listCmd = m.stationsList.Update(msg)
	}

	m.updateStations()

	return m, tea.Batch(listCmd)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	topView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.countriesList.View(),
		m.stationsList.View(),
	)

	return topView
}
