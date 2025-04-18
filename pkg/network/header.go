package network

import (
	"net"
	"strconv"
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

		if validateResponse(conn, httpVersion, allowEmpty) {
			return true
		}
	}
	return false
}

func validateResponse(conn net.Conn, httpVersion string, allowEmpty bool) bool {
	// read and validate response
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Err(err, "Read data from socket.")
			return true
		}
		log.Info("Read header response bytes", n)

		if n > 0 {
			return isValidResponse(data, httpVersion)
		}

		if allowEmpty {
			return true
		}

		time.Sleep(time.Duration(1) * time.Second)
	}

}

func isValidResponse(data []byte, httpVersion string) bool {
	dataString := string(data[:])
	if !strings.HasPrefix(dataString, "HTTP/"+httpVersion) {
		log.Error("Not a valid http response!")
		return false
	}

	split := strings.Split(dataString, "\r\n")
	if len(split) < 2 {
		log.Error("not sufficient components in response.", dataString)
		return false
	}

	code := strings.Split(split[1], " ")
	if len(code) < 2 {
		log.Error("no valid status code in", dataString, "\n", code)
		return false
	}

	if val, err := strconv.Atoi(code[0]); err == nil {
		if val < 300 {
			return true
		}
	}

	// for _, key := range split {
	// 	if key == "100 Continue" {
	// 		return true
	// 	}
	// }
	// log.Error("No 100 Continue!")

	log.Error("invalid response", dataString)
	return false
}
