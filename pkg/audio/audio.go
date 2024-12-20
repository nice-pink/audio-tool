package audio

import (
	"os"

	"github.com/nice-pink/audio-tool/pkg/audio/encodings"
	"github.com/nice-pink/goutil/pkg/log"
)

type AudioFile struct {
	Filepath string
	Data     []byte
	Infos    *encodings.AudioInfos
}

func NewAudioFile(filepath string, loadData, parse bool) *AudioFile {
	audio := &AudioFile{Filepath: filepath}

	if loadData || parse {
		data, err := os.ReadFile(filepath)
		if err != nil {
			log.Err(err, "could not load file")
			return nil
		}
		audio.Data = data
	}

	if parse {
		parser := encodings.NewParser()
		audio.Infos = parser.Parse(audio.Data, filepath, false, false, false)
	}

	return audio
}
