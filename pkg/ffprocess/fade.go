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

func fade(in bool, procJob models.ProcJob, codecConfig CodecConfig) error {
	// get propper codec name
	args := GetKwArgs(procJob.AudioFormat, codecConfig)

	// filter args
	filterKwArgs := ffmpeg.KwArgs{}

	switch procJob.ProcInfo.TimeFormat {
	case models.TimeFormat_Samples:
		filterKwArgs["start_sample"] = strconv.FormatFloat(procJob.ProcInfo.Offset, 'f', 4, 64)
		filterKwArgs["nb_samples"] = strconv.FormatFloat(procJob.ProcInfo.Duration, 'f', 4, 64)
	case models.TimeFormat_Seconds:
		filterKwArgs["start_time"] = strconv.FormatFloat(procJob.ProcInfo.Offset, 'f', 4, 64) + "s"
		filterKwArgs["duration"] = strconv.FormatFloat(procJob.ProcInfo.Duration, 'f', 4, 64) + "s"
	}

	filterKwArgs["silence"] = strconv.FormatFloat(procJob.ProcInfo.From, 'f', 4, 64)
	filterKwArgs["unity"] = strconv.FormatFloat(procJob.ProcInfo.To, 'f', 4, 64)

	// transcode
	err := ffmpeg.Input(procJob.Input).
		Filter("afade", ffmpeg.Args{}, filterKwArgs).
		Output(procJob.Output, args).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Err(err)
	}
	return err
}
