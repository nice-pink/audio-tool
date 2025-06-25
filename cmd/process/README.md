# What?

Use ffmpeg for
- transcode audio file.
- fade in/out audio file.
- mix multiple audio files.

# Examples

## Fade job

Codec      string `json:"codec,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	SampleRate int    `json:"sampleRate,omitempty"`
	Channels   int    `json:"channels,omitempty"`

bin/process -job '{"input":"bin/elefanten.mp3","output":"bin/output.mp3","type":"fadeIn","procInfo":{"offset":0.0,"duration":3.0,"from":0.0,"to":1.0},"audioFormat":{"codec":"mp3","bitrate":128000,"sampleRate":44100,"channels":2}}' -codecConfig cmd/process/codec-config.yaml