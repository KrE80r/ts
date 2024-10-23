package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kre80r/ts/pkg/transcript"
)

func main() {
	var audioFile string
	var model string
	var listModels bool

	// Define CLI flags
	flag.StringVar(&audioFile, "audio_file", "", "Path to the audio file to be transcribed.")
	flag.StringVar(&model, "model", "", "Select the model to use for transcription.")
	flag.BoolVar(&listModels, "listmodels", false, "List all available models.")
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Load(os.ExpandEnv("$HOME/.config/fabric/.env"))
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
		return
	}

	// Initialize Whisper with the chosen model
	whisper, err := transcript.NewWhisper(model)
	if err != nil {
		fmt.Println("Error initializing Whisper:", err)
		return
	}

	// Handle CLI options
	if listModels {
		whisper.ListModels()
	} else if audioFile != "" {
		whisper.ProcessFile(audioFile)
	} else {
		fmt.Println("Please provide an audio file or use --listmodels.")
	}
}
