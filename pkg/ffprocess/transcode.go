package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// transcode

func Transcode(procJob models.ProcJob, codecConfig CodecConfig) error {
	// transcode
	inputNode := ffmpeg.Input(procJob.Input).ASplit()

	// get multiple outputs
	outs := GetOutputs(inputNode, procJob.Outputs, procJob.AudioFormats, codecConfig)

	// run
	err := ffmpeg.MergeOutputs(outs...).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Err(err)
	}
	return err
}
