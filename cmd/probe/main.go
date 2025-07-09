package main

import (
	"flag"
	"os"

	"github.com/nice-pink/audio-tool/pkg/ffprocess"
	"github.com/nice-pink/goutil/pkg/log"
)

// analyse audio files
// - metadata / tags

func main() {
	filepath := flag.String("input", "", "Path to file.")
	frames := flag.Bool("frames", false, "Show frames.")
	verbose := flag.Bool("verbose", false, "Show verbose output.")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	info, err := ffprocess.Probe(*filepath, *frames)
	if err != nil {
		log.Err(err, "Could not probe audio")
		os.Exit(1)
	}
	if *verbose {
		log.Info("Probe audio:")
		log.Info(info)
	}
}
