package stream

import (
	"strconv"
)

type IcyMeta struct {
	Bitrate    int
	SampleRate int
	Channels   int
	Quality    int
	AudioInfo  string
	Url        string
}

func GetIcecastPutHeader(connTarget ConnTarget, meta IcyMeta, httpVersion string) ([]byte, error) {
	header := "PUT " + connTarget.MountPoint + " HTTP/" + httpVersion + "\n" +
		"Host: " + connTarget.Domain + ":" + connTarget.Port + "\n" +
		"User-Agent: " + connTarget.UserAgent + "\n"
	if connTarget.BasicAuth != "" {
		header += "Authorization: Basic " + connTarget.BasicAuth + "\n"
	}
	header += addIcyMeta(meta)
	return convertToByteHeader(header, false), nil
}

func addIcyMeta(meta IcyMeta) string {
	audioInfo := "samplerate=" + strconv.Itoa(meta.SampleRate) + ";quality=" + strconv.Itoa(meta.Quality) + ";channels=" + strconv.Itoa(meta.Channels)
	return `Content-Type: audio/mpeg
Accept: */*
User-Agent: streamey
Server: Icecast 2.4.0-kh15
icy-br:` + strconv.Itoa(meta.Bitrate) + `
icy-genre:Test
icy-name:SineSweep
icy-notice1:This is a radiosphere test stream.
icy-pub:0
icy-url:` + meta.Url + `
Icy-MetaData:0
icy-audio-info:` + audioInfo + `
ice-audio-info:` + audioInfo + `
Expect: 100-continue
`
}
