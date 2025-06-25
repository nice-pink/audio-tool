package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// transcode

func TranscodeAudio(filepathIn string, filepathOut string, format models.AudioFormat, codecConfig CodecConfig) error {
	// get propper codec name
	args := GetKwArgs(format, codecConfig)
	// transcode
	err := ffmpeg.Input(filepathIn).
		Output(filepathOut, args).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Err(err)
	}
	return err
}
