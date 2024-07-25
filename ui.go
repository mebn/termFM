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

func initialModel() model {
	// countries
	items := []list.Item{}
	countries := goradios.FetchCountriesDetailed(goradios.OrderName, false, false)
	for _, country := range countries {
		items = append(items, countryItem{title: country.Name, desc: fmt.Sprint(country.StationCount, " stations")})
	}
	countriesView := list.New(items, list.NewDefaultDelegate(), 0, 0)

	// stations
	stations := make(map[string][]list.Item)

	countriesView.Title = "Countries"
	countriesView.SetShowHelp(false)
	countriesView.Paginator.KeyMap.NextPage.SetEnabled(false)
	countriesView.Paginator.KeyMap.PrevPage.SetEnabled(false)
	countriesView.SetShowPagination(false)

	items = []list.Item{}
	for _, station := range goradios.FetchStations(goradios.StationsByCountry, countriesView.SelectedItem().FilterValue()) {
		items = append(items, stationItem{title: station.Name, desc: fmt.Sprintf("%s, %s", station.Country, station.State), url: station.URLResolved})
	}
	stations[countriesView.SelectedItem().FilterValue()] = items
	stationsView := list.New(stations[countriesView.SelectedItem().FilterValue()], list.NewDefaultDelegate(), 0, 0)

	stationsView.Title = "Stations"
	stationsView.SetShowHelp(false)
	stationsView.Paginator.KeyMap.NextPage.SetEnabled(false)
	stationsView.Paginator.KeyMap.PrevPage.SetEnabled(false)

	return model{
		countriesView: countriesView,
		stations:      stations,
		stationsView:  stationsView,
		playerView:    "prev play/pause next",
		quitting:      false,
		focusedView:   country,
	}
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
			if m.focusedView == station {
				m.focusedView = country
			}

		case "left", "h":
			if m.focusedView == country {
				m.focusedView = station
			}
		}

		// focus on view
		if m.focusedView == country {
			var cmd tea.Cmd
			m.countriesView, cmd = m.countriesView.Update(msg)
			return m, cmd
		} else if m.focusedView == station {
			var cmd tea.Cmd
			m.stationsView, cmd = m.stationsView.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.countriesView.SetSize(msg.Width, msg.Height)
		m.stationsView.SetSize(msg.Width, msg.Height)
	}

	// update stationsView
	val, ok := m.stations[m.countriesView.SelectedItem().FilterValue()]

	if ok {
		m.stationsView.SetItems(val)
	} else {
		items := []list.Item{}
		stations := goradios.FetchStations(goradios.StationsByCountry, m.countriesView.SelectedItem().FilterValue())
		for _, station := range stations {
			items = append(items, stationItem{
				title: station.Name,
				desc:  fmt.Sprintf("%s, %s", station.Country, station.State),
				url:   station.URLResolved,
			})
		}
		m.stations[m.countriesView.SelectedItem().FilterValue()] = items
		m.stationsView.SetItems(val)
	}

	m.countriesView, cmd1 = m.countriesView.Update(msg)
	m.stationsView, cmd2 = m.stationsView.Update(msg)

	return m, tea.Batch(cmd1, cmd2)
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
