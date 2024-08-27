package cli

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/mebn/termfm/internal/audioplayer"
	"gitlab.com/AgentNemo/goradios"
)

const HowToUse string = `Usage: termfm [options]
A TUI or CLI interface for listening to radio in the terminal.
If no flags are present, TUI mode will be displayed.
Flags can be combined, e.g. termfm -nr

Options:
	-h, --help : Display this help text.
	-n : CLI mode, No TUI interface. Without this flag, all other flags does nothing.
	-r : Play a random country and radio station.
`

var FlagNotFoundError = fmt.Errorf("Flag not found.")

type Config struct {
	Cli    bool
	random bool
}

func NewConfig() Config {
	return Config{
		Cli:    false,
		random: false,
	}
}

// Takes the command line arguments, excluding the first entry (os.Args[0]).
//
// This function handles all the possible flags the program can take.
func (c *Config) HandleConfig(args []string) error {
	// handle flags
	for i := range args {
		if args[i] == "--help" {
			showHelp()
		} else if args[i][0] != '-' {
			break
		}

		for _, char := range args[i][1:] {
			switch char {
			case 'h':
				showHelp()
			case 'n':
				c.Cli = true
			case 'r':
				c.random = true
			default:
				return fmt.Errorf("-%c: %w", char, FlagNotFoundError)
			}
		}
	}

	return nil
}

// Print help text and exit with status code 0.
func showHelp() {
	fmt.Fprint(os.Stderr, HowToUse)
	os.Exit(0)
}

func (c *Config) ShowCLI() {
	player := audioplayer.NewPlayer()
	stations := goradios.FetchAllStations()
	fmt.Println("starting loop")
	for {
		i := rand.IntN(len(stations))
		if stations[i].Codec != "MP3" {
			continue
		}
		fmt.Println("New station", stations[i].URLResolved)
		go player.Play(stations[i].URLResolved)
		time.Sleep(time.Second * 5)
	}
}
