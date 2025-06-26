package main

import (
	"flag"
	"io"
	"os"

	"github.com/nice-pink/audio-tool/pkg/audio/encodings"
	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
)

// remove tags from audio files

func main() {
	filepath := flag.String("input", "", "Input filepath")
	output := flag.String("output", "", "Output filepath")
	verbose := flag.Bool("verbose", false, "Make output verbose.")
	flag.Parse()

	// get file data

	if *filepath == "" {
		flag.Usage()
		os.Exit(2)
	}

	file, err := os.Open(*filepath)
	if err != nil {
		log.Err(err, "Cannot open file.")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Err(err, "Cannot read file.")
	}

	// parse audio
	parser := encodings.NewParser()
	info := parser.Parse(data, *filepath, false, *verbose, true)
	if info.TagSize > 0 && *output != "" {
		util.WriteDataToFile(data[info.TagSize:], *output)
	}
}
