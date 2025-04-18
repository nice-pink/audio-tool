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

func GetIcecastPutHeader(connTarget ConnTarget, meta IcyMeta, httpVersion string, print bool) ([]byte, error) {
	header := "PUT " + connTarget.MountPoint + " HTTP/" + httpVersion + "\r\n"
	header += "Host: " + connTarget.Domain + ":" + connTarget.Port + "\r\n"
	header += "User-Agent: " + connTarget.UserAgent + "\r\n"
	if connTarget.BasicAuth != "" {
		header += "Authorization: Basic " + connTarget.BasicAuth + "\r\n"
	}
	header += addIcyMeta(meta)
	header += "\r\n"
	return convertToByteHeader(header, print), nil
}

func addIcyMeta(meta IcyMeta) string {
	audioInfo := "samplerate=" + strconv.Itoa(meta.SampleRate) + ";quality=" + strconv.Itoa(meta.Quality) + ";channels=" + strconv.Itoa(meta.Channels)
	add := "Content-Type: audio/mpeg\r\n"
	add += "Accept: */*\r\n"
	add += "Server: Icecast 2.4.0-kh15\r\n"
	add += "icy-br:" + strconv.Itoa(meta.Bitrate) + "\r\n"
	add += "icy-genre:Test\r\n"
	add += "icy-name:SineSweep\r\n"
	add += "icy-notice1:This is a radiosphere test stream.\r\n"
	add += "icy-pub:0\r\n"
	add += "icy-url:" + meta.Url + "\r\n"
	add += "Icy-MetaData:0\r\n"
	add += "icy-audio-info:" + audioInfo + "\r\n"
	add += "ice-audio-info:" + audioInfo + "\r\n"
	add += "Expect: 100-continue\r\n"
	return add
}
