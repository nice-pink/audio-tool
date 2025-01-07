package audio

import (
	"os"

	"github.com/nice-pink/audio-tool/pkg/audio/encodings"
	"github.com/nice-pink/goutil/pkg/log"
)

type AudioFile struct {
	Filepath string
	Size     int64
	Data     []byte
	Infos    *encodings.AudioInfos
}

func NewAudioFile(filepath string, loadData, parse bool) *AudioFile {
	audio := &AudioFile{Filepath: filepath}

	fs, err := os.Stat(filepath)
	if err != nil {
		log.Err(err, "file stats", filepath)
		return nil
	}
	audio.Size = fs.Size()

	if loadData || parse {
		data, err := os.ReadFile(filepath)
		if err != nil {
			log.Err(err, "could not load file data", filepath)
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
