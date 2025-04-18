package stream

func GetShoutcastSourceHeader(connTarget ConnTarget, httpVersion string) ([]byte, error) {
	header := "SOURCE " + connTarget.MountPoint + " HTTP/" + httpVersion + "\n"
	header += "Host: " + connTarget.Domain + ":" + connTarget.Port + "\n"
	header += "User-Agent: " + connTarget.UserAgent + "\n"
	if connTarget.BasicAuth != "" {
		header += "Authorization: Basic " + connTarget.BasicAuth + "\n"
	}
	header += "Connection: close"
	return convertToByteHeader(header, false), nil
}
