package ffprocess

import (
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// codecs

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
	if formatCodec == "" {
		return "copy"
	}

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

// outputs

func GetAudioFormatArgs(format models.AudioFormat, codecConfig CodecConfig) ffmpeg.KwArgs {
	args := ffmpeg.KwArgs{
		"c:a": GetCodec(format.Codec, codecConfig),
	}
	if format.Channels > 0 {
		args["ac"] = format.Channels
	}
	if format.Bitrate > 0 {
		args["b:a"] = format.Bitrate
	}
	if format.SampleRate > 0 {
		args["ar"] = format.SampleRate
	}

	return args
}

func GetOutputs(input *ffmpeg.Node, outputs []string, formats []models.AudioFormat, codecConfig CodecConfig) []*ffmpeg.Stream {
	if len(outputs) != len(formats) {
		log.Error("amount of outputs must match amount of audio formats")
		return nil
	}

	outs := []*ffmpeg.Stream{}
	for index := range len(outputs) {
		args := GetAudioFormatArgs(formats[index], codecConfig)
		out := input.Get(strconv.Itoa(index)).Output(outputs[index], args)
		outs = append(outs, out)
	}

	return outs
}

func GenerateStreamOutputPipeline(stream *ffmpeg.Stream, outputs []string, formats []models.AudioFormat, index int, codecConfig CodecConfig) *ffmpeg.Stream {
	if len(outputs) != len(formats) {
		log.Error("amount of outputs must match amount of audio formats")
		return nil
	}
	if index >= len(outputs) {
		return stream
	}
	stream = stream.Output(outputs[index], GetAudioFormatArgs(formats[index], codecConfig))
	return GenerateStreamOutputPipeline(stream, outputs, formats, index+1, codecConfig)
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
