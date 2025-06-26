package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// transcode
// -f segment -segment_time 10 -segment_format_options movflags=+faststart out%03d.mp4

func Segment(procJob models.ProcJob, codecConfig CodecConfig) error {
	// transcode
	inputNode := ffmpeg.Input(procJob.Input).ASplit()

	// run
	return RunFFmpegInputNode(inputNode, procJob, codecConfig)
}
