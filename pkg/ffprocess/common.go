package ffprocess

import (
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nice-pink/audio-tool/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/u2takey/go-utils/sets"
)

// consume ffmpeg from pipe: cat file.mp3 | ffmpeg -f mp3 -i pipe: -c:a pcm_s16le -f s16le pipe:

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

func GetAudioOutputArgs(format models.AudioFormat, output models.Output, codecConfig CodecConfig) ffmpeg.KwArgs {
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

	if output.SegmentDuration > 0 {
		args["f"] = "segment"
		args["segment_time"] = output.SegmentDuration
	}

	return args
}

func GetOutputs(input *ffmpeg.Node, outputs []models.Output, codecConfig CodecConfig) []*ffmpeg.Stream {
	log.Info("outputs:", outputs)

	outs := []*ffmpeg.Stream{}
	for index, output := range outputs {
		args := GetAudioOutputArgs(output.Format, output, codecConfig)
		out := input.Get(strconv.Itoa(index)).Output(output.Filename, args)
		outs = append(outs, out)
	}
	return outs
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

// input

func NewInputsNode(name string, streams []*ffmpeg.Stream, kwargs ffmpeg.KwArgs) *ffmpeg.Node {
	return ffmpeg.NewNode(streams,
		name,
		sets.NewString("FilterableStream"),
		"FilterableStream",
		0,
		0,
		nil,
		kwargs,
		"InputNode")
}

// global

func GetGlobalArgs(procJob models.ProcJob) []string {
	var args []string
	if len(procJob.GlobalParams) > 0 {
		args = procJob.GlobalParams
	} else {
		args = []string{}
	}

	//* Tags *//
	// id3
	if procJob.TagProc.DiscardId3 {
		args = append(args, "-id3v2_version", "0")
	}
	// xing
	if procJob.TagProc.DiscardXing {
		args = append(args, "-write_xing", "0")
	}

	if len(args) == 0 {
		return nil
	}
	return args
}

// run

func RunFFmpegInputNode(input *ffmpeg.Node, procJob models.ProcJob, codecConfig CodecConfig) error {
	// get multiple outputs
	outs := GetOutputs(input, procJob.Outputs, codecConfig)
	return RunFFmpegMultipleOut(procJob, outs...)
}

func RunFFmpegInputNodeSimple(input *ffmpeg.Node, procJob models.ProcJob, codecConfig CodecConfig) error {
	return input.Stream("", "").Output(procJob.Outputs[0].Filename).OverWriteOutput().ErrorToStdOut().Run()
}

func RunFFmpegMultipleOut(procJob models.ProcJob, outs ...*ffmpeg.Stream) error {
	out := ffmpeg.MergeOutputs(outs...)
	if args := GetGlobalArgs(procJob); args != nil {
		out = out.GlobalArgs(args...)
	}

	// run
	err := out.OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Err(err)
	}
	return err
}
