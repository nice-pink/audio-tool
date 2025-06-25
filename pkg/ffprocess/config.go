package ffprocess

// codec config
/***
Example:
codecs:
- mp3:libmp3lame
- aac:libfdk_aac
***/

type CodecConfig struct {
	UseDefault bool
	Codecs     []string
}

type ProcessConfig struct {
	CodecConfig CodecConfig
}
