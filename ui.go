package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/AgentNemo/goradios"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type countryItem struct {
	title, desc string
}

func (i countryItem) Title() string       { return i.title }
func (i countryItem) Description() string { return i.desc }
func (i countryItem) FilterValue() string { return i.title }

type model struct {
	selected  map[int]struct{}
	countries list.Model
	stations  map[string]goradios.Station
}

func initialModel() model {
	items := []list.Item{}
	coutries := goradios.FetchCountriesDetailed(goradios.OrderName, false, false)
	for _, country := range coutries {
		items = append(items, countryItem{title: country.Name, desc: fmt.Sprint(country.StationCount, " stations")})
	}

	model := model{
		selected:  make(map[int]struct{}),
		countries: list.New(items, list.NewDefaultDelegate(), 0, 0),
		stations:  make(map[string]goradios.Station),
	}

	model.countries.Title = "countries"

	return model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("TermFM")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.countries.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.countries, cmd = m.countries.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.countries.View())
}
