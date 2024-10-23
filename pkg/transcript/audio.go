package transcript

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

// SplitAndProcessAudio splits an audio file into segments and processes each one
func (w *Whisper) SplitAndProcessAudio(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening audio file:", err)
		return
	}
	defer file.Close()

	// Initialize MP3 decoder
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error decoding MP3 file:", err)
		return
	}

	// Set segment duration (10 minutes)
	segmentDuration := 10 * time.Minute
	bufferSize := int(decoder.SampleRate() * 4) // Assuming 16-bit stereo samples

	// Process the audio in 10-minute segments
	for segmentStart := int64(0); ; segmentStart += int64(segmentDuration.Seconds()) {
		// Create a buffer for the segment
		segment := make([]byte, bufferSize)

		// Read the segment into the buffer
		n, err := io.ReadFull(decoder, segment)
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			fmt.Println("Error reading audio segment:", err)
			return
		}

		// Truncate the buffer to the actual read size
		segment = segment[:n]

		// Process the segment (transcription)
		w.ProcessSegment(segment)
	}
}

// ProcessSegment handles the transcription of a single audio segment
func (w *Whisper) ProcessSegment(segment []byte) {
	// Write segment to a temporary file
	tmpFile, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write the segment to the temp file
	if _, err := tmpFile.Write(segment); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		fmt.Println("Error closing temporary file:", err)
		return
	}

	// Send the segment to the transcription API
	transcription, err := w.transcribeAudio(tmpFile.Name())
	if err != nil {
		fmt.Println("Error transcribing audio segment:", err)
		return
	}

	// Print the transcription result
	fmt.Println("Transcription:", transcription)
}
