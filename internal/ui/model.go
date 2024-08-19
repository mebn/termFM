package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/AgentNemo/goradios"
)

type focusedView int

const (
	countryView focusedView = iota
	stationView
	playerView
)

type model struct {
	countriesList list.Model
	stations      map[string][]list.Item
	stationsList  list.Model
	quitting      bool
	state         focusedView
}

func (m *model) updateStations() tea.Cmd {
	var cmd tea.Cmd
	items, ok := m.stations[m.countriesList.SelectedItem().FilterValue()]

	if ok {
		cmd = m.stationsList.SetItems(items)
	} else {
		items := []list.Item{}
		fetchedStudents := goradios.FetchStations(
			goradios.StationsByCountry,
			m.countriesList.SelectedItem().FilterValue(),
		)

		for _, station := range fetchedStudents {
			var desc string
			if station.State == "" {
				desc = station.Country
			} else {
				desc = fmt.Sprintf("%s, %s", station.Country, station.State)
			}

			items = append(items, stationItem{
				title: station.Name,
				desc:  desc,
				url:   station.URLResolved,
			})
		}

		m.stations[m.countriesList.SelectedItem().FilterValue()] = items
		cmd = m.stationsList.SetItems(items)
	}

	return cmd
}

func InitialModel() model {
	// countries
	items := []list.Item{}
	countries := goradios.FetchCountriesDetailed(goradios.OrderName, false, false)
	for _, country := range countries {
		if country.Name == "" {
			continue
		}

		var descAmount string
		if country.StationCount == 1 {
			descAmount = " station"
		} else {
			descAmount = " stations"
		}

		items = append(items, countryItem{
			title: country.Name,
			desc:  fmt.Sprint(country.StationCount, descAmount),
		})
	}

	countriesList := list.New(items, list.NewDefaultDelegate(), 0, 0)

	countriesList.Title = "Countries"
	countriesList.KeyMap.PrevPage = key.NewBinding()
	countriesList.KeyMap.NextPage = key.NewBinding()

	// stations
	stations := make(map[string][]list.Item)

	stationsView := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	stationsView.Title = "Stations"
	stationsView.KeyMap.PrevPage = key.NewBinding()
	stationsView.KeyMap.NextPage = key.NewBinding()

	model := model{
		countriesList: countriesList,
		stations:      stations,
		stationsList:  stationsView,
		quitting:      false,
		state:         countryView,
	}

	model.updateStations()

	return model
}
