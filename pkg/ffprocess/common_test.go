package ffprocess

import (
	"testing"
)

func TestGetOutputFilepath(t *testing.T) {
	// change codec_type
	codec := "aac"

	// all set
	outputFolder := "out"
	filenameIn := "in.mp3"
	filepath := GetOutputFilepath(filenameIn, codec, outputFolder)
	want := "out/in.aac"
	if filepath != want {
		t.Error("all set:", filepath, "!=", want)
	}

	// no folder
	outputFolder = ""
	filenameIn = "in.mp3"
	filepath = GetOutputFilepath(filenameIn, codec, outputFolder)
	want = "in.aac"
	if filepath != want {
		t.Error("no folder:", filepath, "!=", want)
	}

	// no ext
	outputFolder = ""
	filenameIn = "in"
	filepath = GetOutputFilepath(filenameIn, codec, outputFolder)
	outputFolder = ""
	want = "in.aac"
	if filepath != want {
		t.Error("no ext:", filepath, "!=", want)
	}

	// cloud storage
	outputFolder = "out"
	filenameIn = "gs://bucket/in.original"
	filepath = GetOutputFilepath(filenameIn, codec, outputFolder)
	want = "out/in.aac"
	if filepath != want {
		t.Error("all set:", filepath, "!=", want)
	}
}

func TestGetCodecDefault(t *testing.T) {
	// mp3
	is := "mp3"
	want := "libmp3lame"
	if GetCodecDefault(is) != want {
		t.Error("mp3:", is, "!=", want)
	}

	//aac
	is = "aac"
	want = "libfdk_aac"
	if GetCodecDefault(is) != want {
		t.Error("aac:", is, "!=", want)
	}

	//flac
	is = "flac"
	want = "flac"
	if GetCodecDefault(is) != want {
		t.Error("flac:", is, "!=", want)
	}
}

func TestGetCodec(t *testing.T) {
	codecConfig := CodecConfig{Codecs: []string{"mp3:libmp3lame", "aac:libfdk_aac"}}

	// mp3
	is := "mp3"
	want := "libmp3lame"
	if GetCodec(is, codecConfig) != want {
		t.Error("mp3:", is, "!=", want)
	}

	//aac
	is = "aac"
	want = "libfdk_aac"
	if GetCodec(is, codecConfig) != want {
		t.Error("aac:", is, "!=", want)
	}

	//flac
	is = "flac"
	want = "flac"
	if GetCodec(is, codecConfig) != want {
		t.Error("flac:", is, "!=", want)
	}
}
