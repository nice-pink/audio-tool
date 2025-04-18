package stream

func GetShoutcastSourceHeader(connTarget ConnTarget, httpVersion string, print bool) ([]byte, error) {
	header := "SOURCE " + connTarget.MountPoint + " HTTP/" + httpVersion + "\r\n"
	header += "Host: " + connTarget.Domain + "\r\n"
	header += "User-Agent: " + connTarget.UserAgent + "\r\n"
	if connTarget.BasicAuth != "" {
		header += "Authorization: Basic " + connTarget.BasicAuth + "\r\n"
	}
	header += "Connection: close" + "\r\n"
	header += "\r\n"
	return convertToByteHeader(header, print), nil
}
