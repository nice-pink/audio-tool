package ffprocess

import (
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
	// input args
	intputKwArgs := ffmpeg.KwArgs{}
	if procJob.ProcInfo.Trim {
		if in {
			intputKwArgs["ss"] = GetFloatString(procJob.ProcInfo.Offset)
		} else {
			intputKwArgs["to"] = GetFloatString(procJob.ProcInfo.Offset + procJob.ProcInfo.Duration)
		}
	}

	// filter args
	filterKwArgs := getFadeFilterArgs(procJob.ProcInfo, in)

	inputNode := ffmpeg.Input(procJob.Input, intputKwArgs).
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
	// case models.TimeFormat_Samples:
	// 	filterKwArgs["start_sample"] = GetFloatString(procInfo.Offset)
	// 	filterKwArgs["nb_samples"] = GetFloatString(procInfo.Duration)
	case models.TimeFormat_Seconds:
		if fadeIn && !procInfo.Trim {
			// if trim=true -> the input will be started later in global command
			filterKwArgs["start_time"] = GetFloatString(procInfo.Offset) + "s"
		}
		filterKwArgs["duration"] = GetFloatString(procInfo.Duration) + "s"
	}

	if fadeIn {
		filterKwArgs["silence"] = GetFloatString(procInfo.From)
		filterKwArgs["unity"] = GetFloatString(procInfo.To)
	} else {
		filterKwArgs["silence"] = GetFloatString(procInfo.To)
		filterKwArgs["unity"] = GetFloatString(procInfo.From)
	}

	if !fadeIn {
		filterKwArgs["type"] = "out"
	}

	return filterKwArgs
}
