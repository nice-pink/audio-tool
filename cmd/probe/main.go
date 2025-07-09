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
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	ffprocess.Probe(*filepath, *frames)
}
