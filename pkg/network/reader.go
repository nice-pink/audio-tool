package network

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
)

const (
	TCP_PROTO = "tcp"
)

type DataValidator interface {
	ValidateData(data []byte) error
	// Validate() error
}

type DummyValidator struct {
}

func (v DummyValidator) ValidateData(data []byte) error {
	return nil
}

// func (v DummyValidator) Validate() error {
// 	return nil
// }

// read stream

func (c Connection) ReadStream(outputFilepath string, reconnect bool, dataValidator DataValidator) {
	// early exit
	if c.url == "" {
		log.Newline()
		log.Error("Define url!")
		flag.Usage()
		os.Exit(2)
	}

	// log infos
	if outputFilepath != "" {
		log.Info("Dump data to file:", outputFilepath)
	}

	for {
		var err error
		log.Newline()
		filepath := util.GetFilePath(outputFilepath)
		if c.connectionType == HttpConnection {
			log.Info("Http connection to url", c.url)
			err = c.ReadHttpLineByLine(filepath, "", dataValidator)
		} else {
			log.Info("Socket connection to url", c.url)
			err = c.ReadSocket(filepath, c.timeout, dataValidator)
		}

		if !reconnect && err != nil {
			break
		}
	}
}

func (c Connection) ReadSocket(dumpToFile string, timeout time.Duration, dataValidator DataValidator) error {
	conn, err := c.GetSocketConn()
	if err != nil {
		log.Err(err, "socket reader can't get connection.")
		return err
	}
	defer conn.Close()

	// open file
	writeToFile := false
	var file *os.File = nil
	if dumpToFile != "" {
		file, err = os.Create(dumpToFile)
		if err != nil {
			log.Err(err, "cannot create file", dumpToFile)
			return err
		}
		writeToFile = true
		defer func() {
			if err := file.Close(); err != nil {
				log.Err(err, "could not close file.")
			}
		}()
	}

	// read data
	var bytesRead int
	var bytes int = 0
	for {
		buffer := make([]byte, 1024)
		bytes, err = conn.Read(buffer)
		if err == io.EOF {
			// done receiving
			return nil
		}
		if err != nil {
			log.Err(err, "socket reader can't read.")
			return err
		}
		if bytes == 0 {
			break
		}
		bytesRead += bytes

		if c.VerboseLogs {
			log.Info(bytes, "bytes read")
		}

		// write to file
		if writeToFile {
			file.Write(buffer[0:bytes])
		}

		// validate
		if dataValidator == nil {
			// skip validation
			continue
		}
		validationErr := dataValidator.ValidateData(buffer[0:bytes])
		if validationErr != nil {
			return validationErr
		}
	}

	return err
}

func (c Connection) ReadHttpLineByLine(dumpToFile string, bearerToken string, dataValidator DataValidator) error {
	// request
	// build request
	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		log.Err(err, "request error.")
		return err
	}

	// auth
	if bearerToken != "" {
		var bearer = "Bearer " + bearerToken
		req.Header.Add("Authorization", bearer)
	}

	// request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Err(err, "client error.")
		return err
	}
	defer resp.Body.Close()

	// open file
	writeToFile := false
	var file *os.File = nil
	if dumpToFile != "" {
		file, err = os.Create(dumpToFile)
		writeToFile = true
		defer func() {
			if err := file.Close(); err != nil {
				log.Err(err, "could not close file.")
			}
		}()
	}

	// read data
	var bytesReadCycle int
	var bytesRead uint64
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Err(err, "could not read bytes")
			return err
		}

		bytesReadCycle = len(line)
		if bytesReadCycle == 0 {
			break
		}
		bytesRead += uint64(bytesReadCycle)

		if c.VerboseLogs {
			log.Info(bytesReadCycle, "bytes read")
		}

		// write to file
		if writeToFile {
			file.Write(line)
		}

		// validate
		if dataValidator == nil {
			// skip validation
			continue
		}
		validationErr := dataValidator.ValidateData(line)
		if validationErr != nil {
			return validationErr
		}
	}

	return err
}

////// quick read test

func ReadTestSocket(port int, dataValidator DataValidator) {
	ln, err := net.Listen(TCP_PROTO, ":"+strconv.Itoa(port))
	if err != nil {
		log.Err(err, "tcp listen")
		return
	}

	// Accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleSocketConnection(conn, dataValidator)
	}
}

func ReadTestHttp(port int, dataValidator DataValidator) {
	listener, err := net.Listen(TCP_PROTO, "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Err(err, "Listen error.")
		return
	}
	defer listener.Close()

	log.Info("Server is listening on port", port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Err(err, "listener accept")
			continue
		}

		// Handle client connection in a goroutine
		go handleHttpClient(conn, dataValidator)
	}
}

func handleSocketConnection(conn net.Conn, dataVaidator DataValidator) {
	// Close the connection when we're done
	defer conn.Close()

	// Read incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Err(err, "connection read")
		return
	}
	dataVaidator.ValidateData(buf)
}

func handleHttpClient(conn net.Conn, dataValidator DataValidator) {
	defer conn.Close()

	var readTotal int64 = 0
	readStart := time.Now().UnixNano()
	// Create a buffer to read data into
	buffer := make([]byte, 1024)

	var count int = 0
	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			log.Err(err, "Read error.")
			return
		}
		readTotal += int64(n)

		// validate
		dataValidator.ValidateData(buffer)

		// Process and use the data (here, we'll just print it)
		if count > 20 {
			rate := ((float64)(readTotal * 8)) / ((float64)(time.Now().UnixNano()-readStart) / 1000000000)
			log.Info("Current rate:", rate)
			count = 0
		}

		count++
	}
}
