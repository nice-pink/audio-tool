package ffprocess

import (
	"testing"

	"github.com/radiosphere/audio-services/pkg/util"
	"github.com/radiosphere/audio-services/test/mock"
)

func TestGetAudioInfo(t *testing.T) {
	_, err := GetAudioInfo(mock.ProbeResultMp3, "")
	if err != nil {
		t.Error("could not parse audio info")
	}
}

func TestHasAudioStream(t *testing.T) {
	audioInfo, err := GetAudioInfo(mock.ProbeResultMp3, "meta")
	if err != nil {
		t.Error("could not parse audio info")
	}

	if !util.StrPtrEqual(audioInfo.Meta, "meta") {
		t.Error("Wrong meta")
	}
	if audioInfo.Streams[0].CodecName != "mp3" {
		t.Error("Wrong codec name")
	}
	if audioInfo.Streams[0].CodecLongName != "MP3 (MPEG audio layer 3)" {
		t.Error("Wrong codec name")
	}
	if audioInfo.Streams[0].BitRate != "128000" {
		t.Error("Wrong bitrate")
	}
	if audioInfo.Streams[0].SampleRate != "44100" {
		t.Error("Wrong sample rate")
	}
	if audioInfo.Streams[0].SampleFormat != "fltp" {
		t.Error("Wrong sample rate")
	}
	if audioInfo.Streams[0].Channels != 2 {
		t.Error("Wrong sample rate")
	}
	if audioInfo.Streams[0].ChannelLayout != "stereo" {
		t.Error("Wrong sample rate")
	}

	// has audio stream
	isAudio := HasAudioStream(audioInfo)
	if !isAudio {
		t.Error("audio stream not identified")
	}

	// change codec_type
	audioInfo.Streams[0].CodecType = "video"
	isNotAudio := HasAudioStream(audioInfo)
	if isNotAudio {
		t.Error("video stream identified as audio")
	}
}
