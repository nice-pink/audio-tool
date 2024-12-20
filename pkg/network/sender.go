package network

import (
	"bufio"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

type Sender struct {
	Url       string
	Port      int
	Meta      []byte
	Filepath  string
	ProxyUrl  string
	ProxyPort int
}

func (s *Sender) SendData() ([]byte, error) {
	addr := s.Url + ":" + strconv.Itoa(s.Port)
	req, err := http.NewRequest(http.MethodPost, addr, nil)
	if err != nil {
		log.Err(err, "create job request")
		return nil, err
	}

	client := s.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err, "post job request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error("status code != 200:", resp.StatusCode, resp.Status)
		return nil, errors.New("status != 200")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err, "read body")
		return nil, err
	}

	return data, nil
}

func (s *Sender) GetHttpClient() *http.Client {
	var client *http.Client
	if s.ProxyUrl != "" && s.ProxyPort != 0 {
		proxyUrl, err := url.Parse("http://" + s.ProxyUrl + ":" + strconv.Itoa(s.ProxyPort))
		if err != nil {
			log.Err(err, "proxy url")
			return nil
		}
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		}
	} else {
		client = &http.Client{}
	}
	return client
}

// helper

func StreamBuffer(address string, sendBitRate float64, byteSegmentSize, fullSize, loops int, buffer []byte) error {
	log.Info("Stream data to", address, "with bitrate", sendBitRate)

	// connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err, "can't dial.")
		return err
	}
	defer conn.Close()

	// variables
	var bytesWrittenCycle int = 0
	var bytesWrittenTotal int = 0
	streamStart := time.Now().UnixNano()
	var byteIndex int = 0
	// var byteSegmentSize int64 = 1024
	bufferLen := fullSize
	if bufferLen <= 0 {
		bufferLen = len(buffer)
	}
	loopCount := 0

	// run loop
	var max int
	var dist int
	var count int = 1
	for {
		if byteIndex >= bufferLen {
			// log.Info("Start loop", loopCount)
			byteIndex = 0
			count = 1
			loopCount++
			if loops > 0 && loopCount >= loops {
				break
			}
		}

		var rate float64 = 0
		if sendBitRate > 0 {
			/*
			* calculate our instant rate over the entire transmit
			* duration
			 */
			rate = ((float64)(bytesWrittenTotal * 8)) / ((float64)(time.Now().UnixNano()-streamStart) / 1000000000)
		}

		// compare rate
		if rate < sendBitRate {
			max = min(bufferLen, count*int(byteSegmentSize))
			dist = max - byteIndex
			// send data
			bytesWrittenCycle, err = conn.Write(buffer[byteIndex:max])
			if err != nil {
				log.Error(err, "could not send data.")
				return err
			}
			if bytesWrittenCycle <= 0 {
				log.Error("bytes written in cycle", bytesWrittenCycle)
				return errors.New("not all data sent")
			}
			if bytesWrittenCycle != dist {
				log.Error("not all bytes sent. Should", dist, ", did", bytesWrittenCycle)
				return errors.New("not all data sent")
			}
			bytesWrittenTotal += bytesWrittenCycle
			byteIndex += bytesWrittenCycle

			count++
		}
	}

	// final log
	streamStop := time.Now().UnixNano()
	passed := streamStart - streamStop
	log.Info("Stopped sending. Bytes:", bytesWrittenTotal, ". Seconds:", passed)
	return nil
}

func SendFile(filepath, address string, byteSegmentSize int) error {
	log.Info("Send file", filepath, "to", address)

	file, err := os.Open(filepath)
	if err != nil {
		log.Error(err, "could not open file", filepath)
		return err
	}
	fs, err := file.Stat()
	if err != nil {
		log.Error(err, "could not get file stats", filepath)
		return err
	}
	fileSize := fs.Size()

	reader := bufio.NewReader(file)

	// connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err, "can't dial.")
		return err
	}
	defer conn.Close()

	// variables
	streamStart := time.Now().UnixNano()
	bytesRead := 0
	bytesWrittenCycle := 0
	bytesWrittenTotal := 0

	var min int
	var bytes int
	// run loop
	for {
		min = int(math.Min(float64(byteSegmentSize), float64(fileSize-int64(bytesRead))))
		if min == 0 {
			break
		}

		buffer := make([]byte, min)
		bytes, err = reader.Read(buffer)
		if err != nil {
			log.Error(err, "could not read data.")
			return err
		}
		bytesRead += bytes

		// send data
		bytesWrittenCycle, err = conn.Write(buffer)
		if err != nil {
			log.Error(err, "could not send data.")
			return err
		}
		bytesWrittenTotal += bytesWrittenCycle
	}

	// final log
	streamStop := time.Now().UnixNano()
	passed := streamStart - streamStop
	log.Info("Stopped sending. Bytes:", bytesWrittenTotal, ". Seconds:", passed)
	return nil
}
