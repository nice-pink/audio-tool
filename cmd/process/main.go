package main

import (
	"flag"

	"github.com/nice-pink/audio-tool/pkg/ffprocess"
)

// fade out:
// ffmpeg -i <input.mp3> -af atrim=end_sample=<end_sample>,afade=type=out:start_sample=<start_sample>:nb_samples=<dur_samples> -b:a <br> -write_xing 0 -id3v2_version 0 <output.mp3>

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

	ffprocess.RunJob(*job, codecConfig)

}
