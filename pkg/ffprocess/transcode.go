package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// transcode

func Transcode(procJob models.ProcJob, codecConfig CodecConfig) error {
	// transcode
	inputNode := ffmpeg.Input(procJob.Input).ASplit()

	// run
	return RunFFmpegInputNode(inputNode, procJob, codecConfig)
}
