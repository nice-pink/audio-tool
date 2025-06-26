package main

import (
	"flag"
	"os"

	"github.com/nice-pink/audio-tool/pkg/ffprocess"
)

func main() {
	// port := flag.Int("port", 8080, "Http port")
	// flag.Parse()
	job := flag.String("job", "", "Job defined in json")
	codecConfigPath := flag.String("codecConfig", "", "Path to codec config.")
	flag.Parse()

	codecConfig := ffprocess.CodecConfig{}
	if *codecConfigPath != "" {

	} else {
		codecConfig.UseDefault = true
	}

	err := ffprocess.RunJob(*job, codecConfig)
	if err != nil {
		os.Exit(2)
	}
}
