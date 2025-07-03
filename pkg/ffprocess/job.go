package ffprocess

import (
	"errors"

	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/data"
	"github.com/nice-pink/goutil/pkg/log"
)

func RunJob(job string, codecConfig CodecConfig) error {
	// get job
	jobMap, err := data.GetJsonMap(job)
	if err != nil {
		log.Err(err, "parse job error")
		return err
	}

	jobType, ok := jobMap["type"]
	if !ok {
		log.Error("job does not have a type")
		return errors.New("invalid job")
	}

	switch jobType {
	// mix job
	case "mix":
		return RunMixJob(job, codecConfig)
	// proc job
	default:
		return RunProcJob(job, codecConfig)
	}
}

func RunProcJob(job string, codecConfig CodecConfig) error {
	var procJob models.ProcJob
	err := data.GetJson(job, &procJob)
	if err != nil {
		log.Err(err, "proc job parsing error")
		return err
	}

	switch procJob.Type {
	// transcode
	case "transcode":
		return Transcode(procJob, codecConfig)
		// fade
	case "fadeIn":
		fallthrough
	case "fadeOut":
		fallthrough
	case "fade":
		return Fade(procJob, codecConfig)
	default:
		log.Error("unknown job:", procJob.Type)
	}

	return nil
}

func RunMixJob(job string, codecConfig CodecConfig) error {

	var mixJob models.MixJob
	err := data.GetJson(job, &mixJob)
	if err != nil {
		log.Err(err, "mix job parsing error")
		return err
	}
	return Mix(mixJob, codecConfig)
}
