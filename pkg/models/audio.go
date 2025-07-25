package models

import (
	"strconv"
)

//

type AudioFormat struct {
	Codec      string `json:"codec,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	SampleRate int    `json:"sampleRate,omitempty"`
	Channels   int    `json:"channels,omitempty"`
}

func GetAudioFormatFromInfo(info AudioInfo) AudioFormat {
	for _, stream := range info.Streams {
		if stream.CodecType == "audio" {
			bitrate, _ := strconv.Atoi(stream.BitRate)
			sampleRate, _ := strconv.Atoi(stream.SampleRate)
			return AudioFormat{
				Codec:      stream.CodecName,
				Bitrate:    bitrate,
				SampleRate: sampleRate,
				Channels:   stream.Channels,
			}
		}
	}
	return AudioFormat{Codec: "", Bitrate: 0, SampleRate: 0, Channels: 0}
}

// proc

type TimeFormat int

const (
	TimeFormat_Seconds TimeFormat = iota
)

type ProcInfo struct {
	Offset     float64
	Duration   float64
	From       float64
	To         float64
	TimeFormat TimeFormat // always uses seconds
	Trim       bool
}

type ProcJob struct {
	Type         string
	Input        string
	ProcInfo     ProcInfo
	Outputs      []Output
	TagProc      TagProc
	GlobalParams []string
}

type Input struct {
	Filename string
	Offset   float64
}

type Output struct {
	Filename        string
	SegmentDuration float64
	Format          AudioFormat
}

type TagProc struct {
	DiscardId3  bool
	DiscardXing bool
}

type MixJob struct {
	Type         string
	Inputs       []Input
	ProcInfos    []ProcInfo
	ProcJob      ProcJob
	GlobalParams []string
}

// ffmpeg info audio info (ffprobe)

type AudioInfo struct {
	Frames  []Frame      `json:"frames,omitempty"`
	Streams []Stream     `json:"streams,omitempty"`
	Format  StreamFormat `json:"format,omitempty"`
	Meta    *string      `json:"meta,omitempty"`
}

func (info AudioInfo) IsValid() bool {
	duration, _ := strconv.ParseFloat(info.Format.DurationSec, 64)
	if len(info.Streams) == 0 {
		return false
	}
	if duration == 0 {
		return false
	}
	return true
}

type StreamFormat struct {
	DurationSec string `json:"duration"`
}

// from ffprobe
type Stream struct {
	Index         int    `json:"index"`
	CodecType     string `json:"codec_type"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	SampleFormat  string `json:"sample_fmt"`
	SampleRate    string `json:"sample_rate"`
	Channels      int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
	Duration      string `json:"duration"`
	BitRate       string `json:"bit_rate"`
}

type Frame struct {
	MediaType               string `json:"media_type"`
	StreamIndex             int    `json:"stream_index"`
	KeyFrame                int    `json:"key_frame"`
	PTS                     int64  `json:"pts"`
	PTSTime                 string `json:"pts_time"`
	PKTDTS                  int64  `json:"pkt_dts"`
	PKTDTSTime              string `json:"pkt_dts_time"`
	BestEffortTimestamp     int64  `json:"best_effort_timestamp"`
	BestEffortTimestampTime string `json:"best_effort_timestamp_time"`
	Duration                int64  `json:"duration"`
	DurationTime            string `json:"duration_time"`
	PKTPos                  string `json:"pkt_pos"`
	PKTSize                 string `json:"pkt_size"`
	SampleFmt               string `json:"sample_fmt"`
	NbSamples               int    `json:"nb_samples"`
	Channels                int    `json:"channels"`
	ChannelLayout           string `json:"channel_layout"`
}

// process job

type AudioJobSpec struct {
	Input   string
	Output  string `json:"output,omitempty"`
	Bitrate int    `json:"bitrate"`
}

// transcribe

type Transcription struct {
	Filepath     string
	Text         string
	Confidence   float32
	BytesWritten int
}

func (t *Transcription) GetConfidencePercentage() int {
	return int(t.Confidence * 100)
}

//

type AudioMetadata struct {
	Filepath string
	Title    string
	Channel  string
	AudioId  string
	Source   string
}

func (a *AudioMetadata) IsValid() bool {
	if a.Filepath == "" || a.Title == "" || a.Channel == "" {
		return false
	}
	return true
}

func (a *AudioMetadata) String() string {
	return "title: " + a.Title + ", file: " + a.Filepath + ", audioId: " + a.AudioId
}
