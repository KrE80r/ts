package transcript

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Whisper struct for handling transcription
type Whisper struct {
	Model     string
	OpenAIKey string
	Client    *http.Client
}

// NewWhisper initializes the Whisper struct with the given model
func NewWhisper(model string) (*Whisper, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return nil, fmt.Errorf("No OpenAI API key found in environment variables")
	}

	if model == "" {
		model = "whisper-1" // Set default model if not specified
	}

	return &Whisper{
		Model:     model,
		OpenAIKey: openaiKey,
		Client:    &http.Client{},
	}, nil
}

// ListModels prints available models for transcription
func (w *Whisper) ListModels() {
	fmt.Println("Available models for transcription:")
	fmt.Println(" - whisper-1")
	// Add more models as needed
}

// ProcessFile handles splitting and processing the audio file for transcription
func (w *Whisper) ProcessFile(filePath string) {
	w.SplitAndProcessAudio(filePath)
}

// transcribeAudio sends the audio file to the OpenAI API for transcription
func (w *Whisper) transcribeAudio(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file for transcription: %v", err)
	}
	defer file.Close()

	// Prepare request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", file)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.OpenAIKey)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Model", w.Model)

	// Send request
	resp, err := w.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed transcription, status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}
