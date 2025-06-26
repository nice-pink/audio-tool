package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Mix(procJob models.ProcJob, codecConfig CodecConfig) error {
	// args := getMixFilterArgs()
	return fade(true, procJob, codecConfig)
}

func getMixFilterArgs(procInfo models.ProcInfo, fadeIn bool) ffmpeg.KwArgs {
	filterKwArgs := ffmpeg.KwArgs{}

	return filterKwArgs
}
