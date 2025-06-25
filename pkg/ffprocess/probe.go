package ffprocess

import (
	"encoding/json"
	"strings"

	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Probe(filepathIn string) (string, error) {
	info, err := ffmpeg.Probe(filepathIn)
	if err != nil {
		log.Err(err, "Could not probe audio")
		return "", err
	}
	log.Info("Probe audio:")
	log.Info(info)
	return info, err
}

func GetAudioInfo(probeInfo, meta string) (*models.AudioInfo, error) {
	var info models.AudioInfo
	if err := json.Unmarshal([]byte(probeInfo), &info); err != nil {
		log.Error(probeInfo)
		log.Err(err, "Cannot unmarshal json")
		return nil, err
	}
	info.Meta = &meta
	return &info, nil
}

func HasAudioStream(audioInfo *models.AudioInfo) bool {
	for _, stream := range audioInfo.Streams {
		if strings.EqualFold(stream.CodecType, "audio") {
			return true
		}
	}
	return false
}
