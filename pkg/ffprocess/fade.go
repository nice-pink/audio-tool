package ffprocess

import (
	"strconv"

	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func FadeIn(procJob models.ProcJob, codecConfig CodecConfig) error {
	return fade(true, procJob, codecConfig)
}

func FadeOut(procJob models.ProcJob, codecConfig CodecConfig) error {
	return fade(false, procJob, codecConfig)
}

func fade(in bool, procJob models.ProcJob, codecConfig CodecConfig) error {
	// filter args
	filterKwArgs := getFadeFilterArgs(procJob.ProcInfo, in)

	inputNode := ffmpeg.Input(procJob.Input).
		Filter("afade", ffmpeg.Args{}, filterKwArgs).ASplit()

	// get multiple outputs
	outs := GetOutputs(inputNode, procJob.Outputs, procJob.AudioFormats, codecConfig)

	// run
	err := ffmpeg.MergeOutputs(outs...).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Err(err)
	}
	return err
}

func getFadeFilterArgs(procInfo models.ProcInfo, fadeIn bool) ffmpeg.KwArgs {
	filterKwArgs := ffmpeg.KwArgs{}

	switch procInfo.TimeFormat {
	case models.TimeFormat_Samples:
		filterKwArgs["start_sample"] = strconv.FormatFloat(procInfo.Offset, 'f', 4, 64)
		filterKwArgs["nb_samples"] = strconv.FormatFloat(procInfo.Duration, 'f', 4, 64)
	case models.TimeFormat_Seconds:
		filterKwArgs["start_time"] = strconv.FormatFloat(procInfo.Offset, 'f', 4, 64) + "s"
		filterKwArgs["duration"] = strconv.FormatFloat(procInfo.Duration, 'f', 4, 64) + "s"
	}

	filterKwArgs["silence"] = strconv.FormatFloat(procInfo.From, 'f', 4, 64)
	filterKwArgs["unity"] = strconv.FormatFloat(procInfo.To, 'f', 4, 64)

	if !fadeIn {
		filterKwArgs["type"] = "out"
	}

	return filterKwArgs
}
