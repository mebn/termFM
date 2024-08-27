package audioplayer

import (
	"net/http"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type AudioPlayer struct {
	otoContext *oto.Context
	player     *oto.Player
}

func NewPlayer() AudioPlayer {
	a := AudioPlayer{
		player: nil,
	}
	a.initOtoContext()

	return a
}

func (a *AudioPlayer) Play(url string) {
	a.Stop()

	resp, err := http.Get(url)
	if err != nil {
		// return fmt.Errorf("Error fetching the live audio stream: %w", err)
		// fmt.Println("Error fetching the live audio stream: %w", err)
		return
	}
	defer resp.Body.Close()

	decoder, err := mp3.NewDecoder(resp.Body)
	if err != nil {
		// return fmt.Errorf("Error decoding the audio stream: %w", err)
		// fmt.Println("Error decoding the audio stream: %w", err)
		return
	}

	player := a.otoContext.NewPlayer(decoder)
	a.player = player

	a.player.Play()

	for a.player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

func (a *AudioPlayer) Stop() {
	if a.player != nil {
		a.player.Close()
	}
}

func (a *AudioPlayer) initOtoContext() {
	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100
	// op.SampleRate = 48000

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoContext, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}

	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	a.otoContext = otoContext
}
