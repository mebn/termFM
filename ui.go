package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/AgentNemo/goradios"
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

type focusedView int

const (
	country focusedView = iota
	station
	player
)

type model struct {
	countriesView list.Model
	stations      map[string][]list.Item
	stationsView  list.Model
	playerView    string
	quitting      bool
	focusedView   focusedView
}

func (m *model) updateStations() tea.Cmd {
	items := []list.Item{}
	fetchedStudents := goradios.FetchStations(goradios.StationsByCountry, m.countriesView.SelectedItem().FilterValue())
	for _, station := range fetchedStudents {
		items = append(items, stationItem{
			title: station.Name,
			desc:  fmt.Sprintf("%s, %s", station.Country, station.State),
			url:   station.URLResolved,
		})
	}
	m.stations[m.countriesView.SelectedItem().FilterValue()] = items
	cmd := m.stationsView.SetItems(items)

	return cmd
}

func initialModel() model {
	// countries
	items := []list.Item{}
	countries := goradios.FetchCountriesDetailed(goradios.OrderName, false, false)
	for _, country := range countries {
		items = append(items, countryItem{title: country.Name, desc: fmt.Sprint(country.StationCount, " stations")})
	}

	countriesView := list.New(items, list.NewDefaultDelegate(), 0, 0)

	countriesView.Title = "Countries"
	countriesView.SetShowHelp(false)

	// stations
	stations := make(map[string][]list.Item)

	stations["Sweden"] = []list.Item{
		stationItem{title: "marcus", desc: "asd", url: "asdasd"},
	}

	stationsView := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	stationsView.Title = "Stations"
	stationsView.SetShowHelp(false)

	model := model{
		countriesView: countriesView,
		stations:      stations,
		stationsView:  stationsView,
		playerView:    "prev play/pause next",
		quitting:      false,
		focusedView:   country,
	}

	model.updateStations()

	return model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("TermFM")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd1, cmd2 tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "right", "l":
			if m.focusedView == country {
				m.focusedView = station
			} else {
				m.focusedView = country
			}

		case "left", "h":
			if m.focusedView == country {
				m.focusedView = station
			} else {
				m.focusedView = country
			}
		}

		// focus on view
		var cmd tea.Cmd
		if m.focusedView == country {
			m.countriesView, cmd = m.countriesView.Update(msg)
			return m, cmd
		} else if m.focusedView == station {
			m.stationsView, cmd = m.stationsView.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.countriesView.SetSize(msg.Width, msg.Height)
		m.stationsView.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	val, ok := m.stations[m.countriesView.SelectedItem().FilterValue()]
	if ok {
		cmd = m.stationsView.SetItems(val)
	} else {
		cmd = m.updateStations()
	}

	m.countriesView, cmd1 = m.countriesView.Update(msg)
	m.stationsView, cmd2 = m.stationsView.Update(msg)

	return m, tea.Batch(cmd1, cmd2, cmd)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	topView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.countriesView.View(),
		m.stationsView.View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, topView, m.playerView)
}
