package ffprocess

import (
	"strconv"

	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Mix(procJob models.MixJob, codecConfig CodecConfig) error {
	// args := getMixFilterArgs()
	return mix(procJob, codecConfig)
}

func mix(job models.MixJob, codecConfig CodecConfig) error {
	// input args
	// intputKwArgs := ffmpeg.KwArgs{}

	// // in trim
	// inProcInfo := job.ProcInfos[0]
	// if inProcInfo.Trim {
	// 	intputKwArgs["ss"] = GetFloatString(inProcInfo.Offset)

	// }

	// // out trim
	// outProcInfo := job.ProcInfos[1]
	// if outProcInfo.Trim {
	// 	intputKwArgs["to"] = GetFloatString(outProcInfo.Offset + outProcInfo.Duration)
	// }

	// don't stretch or skew
	// https://superuser.com/questions/1619992/audio-out-of-sync-when-using-ffmpeg-adelay-and-amix
	// Note: global -async is deprecated!
	job.ProcJob.GlobalParams = append(job.ProcJob.GlobalParams, "-async", "1")

	streams := []*ffmpeg.Stream{}
	for i, in := range job.Inputs {
		procInfos := job.ProcInfos[i]

		// input
		s := ffmpeg.Input(in.Filename, nil)

		if procInfos.Duration > 0 && procInfos.From != procInfos.To {
			s = s.Filter("afade", nil, getFadeFilterArgs(procInfos))
		}

		if in.Offset > 0 {
			// only apply delay filter if necessary
			s = s.Filter("adelay", nil, getDelayFilterArgs(in.Offset))
		}
		streams = append(streams, s)
	}

	filterNodes := ffmpeg.Filter(streams, "amix", nil, getMixFilterArgs()).ASplit()

	// run
	return RunFFmpegInputNode(filterNodes, job.ProcJob, codecConfig)
}

func getDelayFilterArgs(val float64) ffmpeg.KwArgs {
	filterKwArgs := ffmpeg.KwArgs{}

	filterKwArgs["delays"] = strconv.Itoa(GetMilliSeconds(val))
	// delay all channels equally
	filterKwArgs["all"] = "1"

	return filterKwArgs
}

func getMixFilterArgs() ffmpeg.KwArgs {
	filterKwArgs := ffmpeg.KwArgs{}
	return filterKwArgs
}
