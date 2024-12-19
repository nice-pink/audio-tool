package main

import (
	"flag"
	"os"

	"github.com/nice-pink/audio-tool/pkg/tags/id3"
	"github.com/nice-pink/goutil/pkg/log"
)

func main() {
	filepath := flag.String("input", "", "Path to file.")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	data, err := os.ReadFile(*filepath)
	if err != nil {
		log.Err(err, "open file", *filepath)
		os.Exit(2)
	}

	if id3.HasTagId(data) {
		log.Info("has id3 tag")
		size := id3.GetTagSize(data)
		log.Info("tag size:", size, "bytes")
	}
}
