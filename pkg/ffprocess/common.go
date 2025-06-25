package ffprocess

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/nice-pink/audio-tool/pkg/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetKwArgs(format models.AudioFormat, codecConfig CodecConfig) ffmpeg.KwArgs {
	args := ffmpeg.KwArgs{
		"c:a": GetCodec(format.Codec, codecConfig),
		"b:a": format.Bitrate,
		"ar":  format.SampleRate,
	}
	if format.Channels > 0 {
		args["ac"] = format.Channels
	}

	return args
}

func GetCodecDefault(formatCodec string) string {
	if strings.EqualFold(formatCodec, "mp3") {
		return "libmp3lame"
	}
	if strings.EqualFold(formatCodec, "aac") {
		return "libfdk_aac"
	}
	return formatCodec
}

func GetCodec(formatCodec string, codecConfig CodecConfig) string {
	// get codec from config or return as is
	for _, c := range codecConfig.Codecs {
		ref := strings.Split(c, ":")
		if len(ref) < 2 {
			continue
		}
		if strings.EqualFold(ref[0], formatCodec) {
			return ref[1]
		}
	}
	return formatCodec
}

func GetOutputFilepath(filepathIn string, codec string, outputFolder string) string {
	filename := filepath.Base(filepathIn)
	ext := filepath.Ext(filename)
	filepath := strings.TrimSuffix(filename, ext) + "." + codec
	if outputFolder != "" {
		filepath = path.Join(outputFolder, filepath)
	}
	return filepath
}
