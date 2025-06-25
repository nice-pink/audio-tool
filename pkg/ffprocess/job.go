package ffprocess

import (
	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/data"
	"github.com/nice-pink/goutil/pkg/log"
)

func RunJob(job string, codecConfig CodecConfig) error {
	// get job
	var procJob models.ProcJob
	err := data.GetJson(job, &procJob)
	if err != nil {
		log.Err(err, "parse job error")
		return err
	}

	switch procJob.Type {
	case "fadeIn":
		return FadeIn(procJob, codecConfig)
	default:
		log.Error("unknown job:", procJob.Type)
	}

	return nil
}
