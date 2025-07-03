package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Fade(procJob models.ProcJob, codecConfig CodecConfig) error {
	return fade(procJob, codecConfig)
}

func fade(procJob models.ProcJob, codecConfig CodecConfig) error {
	// input args
	intputKwArgs := ffmpeg.KwArgs{}
	if procJob.ProcInfo.Trim {
		if isFadeIn(procJob.ProcInfo) {
			intputKwArgs["ss"] = GetFloatString(procJob.ProcInfo.Offset)
		} else {
			intputKwArgs["to"] = GetFloatString(procJob.ProcInfo.Offset + procJob.ProcInfo.Duration)
		}
	}

	// filter args
	filterKwArgs := getFadeFilterArgs(procJob.ProcInfo)

	inputNode := ffmpeg.Input(procJob.Input, intputKwArgs).
		Filter("afade", ffmpeg.Args{}, filterKwArgs).ASplit()

	// run
	return RunFFmpegInputNode(inputNode, procJob, codecConfig)
}

func getFadeFilterArgs(procInfo models.ProcInfo) ffmpeg.KwArgs {
	filterKwArgs := ffmpeg.KwArgs{}

	fadeIn := isFadeIn(procInfo)

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

func isFadeIn(procInfo models.ProcInfo) bool {
	return procInfo.From < procInfo.To
}
