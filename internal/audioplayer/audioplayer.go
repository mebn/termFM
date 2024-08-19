package audioplayer

import (
	"io"
	"net/http"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

// TODO: fix non-mp3 streams
func PlayStation(url string) {
	// Fetch the live audio stream
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println("Error fetching the live audio stream:", err)
		return
	}
	defer resp.Body.Close()

	// Decode the MP3 stream
	decoder, err := mp3.NewDecoder(resp.Body)
	if err != nil {
		// fmt.Println("Error decoding the audio stream:", err)
		return
	}

	// Initialize the audio context
	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		// fmt.Println("Error initializing audio context:", err)
		return
	}
	defer context.Close()

	// Create a player
	player := context.NewPlayer()
	defer player.Close()

	// Stream the audio to the player
	buffer := make([]byte, 8192)
	for {
		n, err := decoder.Read(buffer)
		if err != nil && err != io.EOF {
			// fmt.Println("Error reading from audio stream:", err)
			return
		}
		if n == 0 {
			break
		}

		if _, err := player.Write(buffer[:n]); err != nil {
			// fmt.Println("Error playing audio:", err)
			return
		}
	}
}
