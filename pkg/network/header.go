package network

import (
	"net"
	"strings"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

func WriteHeader(conn net.Conn, headerBuffer []byte, retry int, httpVersion string, validate, allowEmpty bool) bool {
	for counter := 0; counter < retry; counter++ {
		n, err := conn.Write(headerBuffer)
		if err != nil {
			log.Err(err, "Could not send data.")
			return false
		}

		if n < len(headerBuffer) {
			log.Error("Did not send entire header.")
			return false
		}
		log.Info("Header written", n)

		if !validate {
			return true
		}

		isValid := validateResponse(conn, httpVersion, allowEmpty)
		if isValid {
			return true
		}
	}
	return false
}

func validateResponse(conn net.Conn, httpVersion string, allowEmpty bool) bool {
	// read and validate response
	var data []byte
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Err(err, "Read data from socket.")
			return true
		}

		if n > 0 {
			return isValidResponse(data, httpVersion)
		}

		log.Info("Read header response bytes", n)
		time.Sleep(time.Duration(2) * time.Second)
	}

}

func isValidResponse(data []byte, httpVersion string) bool {
	dataString := string(data[:])
	if !strings.HasPrefix(dataString, "HTTP/"+httpVersion) {
		log.Error("Not a valid http response!")
		return false
	}

	split := strings.Split(dataString, "\r\n")
	if len(split) <= 0 {
		log.Error("No components in response.")
		return false
	}

	for _, key := range split {
		if key == "100 Continue" {
			return true
		}
	}

	log.Error("No 100 Continue!")
	return false
}
