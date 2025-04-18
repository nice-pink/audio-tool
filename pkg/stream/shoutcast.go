package stream

func GetShoutcastSourceHeader(connTarget ConnTarget, meta IcyMeta, httpVersion string) ([]byte, error) {
	header := "SOURCE " + connTarget.MountPoint + " HTTP/" + httpVersion + "\n" +
		"Host: " + connTarget.Domain + ":" + connTarget.Port + "\n" +
		"User-Agent: " + connTarget.UserAgent + "\n" +
		"Connection: close"
	if connTarget.BasicAuth != "" {
		header += "Authorization: Basic " + connTarget.BasicAuth + "\n"
	}
	return convertToByteHeader(header, false), nil
}
