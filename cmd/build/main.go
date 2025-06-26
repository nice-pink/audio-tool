package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/nice-pink/audio-tool/pkg/audio/tags/id3"
	"github.com/nice-pink/goutil/pkg/log"
)

// add metadata tags to file

func main() {
	input := flag.String("input", "", "Input filepath.")
	output := flag.String("output", "", "Output filepath.")
	tags := flag.String("tags", "", "Comma separated list of tag types and sizes. E.g. id3:1024,riff:60")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	inputData, err := os.ReadFile(*input)
	if err != nil {
		log.Err(err, "open file", *input)
		os.Exit(2)
	}

	outFile, err := os.Create(*output)
	if err != nil {
		log.Err(err, "create file", *output)
		os.Exit(2)
	}

	writeTags(outFile, *tags)

	n, err := outFile.Write(inputData)
	if err != nil {
		log.Err(err, "write file", *output)
		os.Exit(2)
	}
	log.Info(n, "bytes written to", *output)
}

func writeTags(file *os.File, tags string) {
	fs, err := file.Stat()
	if err != nil {
		log.Err(err, "file stats", file.Name())
		return
	}

	// add tags
	for _, t := range strings.Split(tags, ",") {
		log.Info("handle", t)
		var data []byte
		tagInfo := strings.Split(t, ":")
		if strings.EqualFold(tagInfo[0], "id3") {
			size, _ := strconv.Atoi(tagInfo[1])
			data = id3.Build(uint32(size), uint32(fs.Size()))
		}

		// write tag to file
		if len(data) > 0 {
			n, err := file.Write(data)
			if err != nil {
				log.Err(err, "could not write", t)
				continue
			}
			log.Info(n, "bytes written", t)
		}
	}
}
