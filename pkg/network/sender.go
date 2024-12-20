package network

import (
	"bufio"
	"errors"
	"math"
	"os"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

func (c Connection) StreamBuffer(buffer []byte, sendBitRate float64, chunkSize int) error {
	// if sendBitRate == 0, then send as quick as possible

	addr := c.GetAddr()
	log.Info("Stream data to", addr, "with bitrate", sendBitRate)

	// connection
	conn, err := c.GetSocketConn()
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
	bufferLen := len(buffer)
	// loopCount := 0

	// run loop
	var max int
	var dist int
	var count int = 1
	for {
		if byteIndex >= bufferLen {
			// log.Info("Start loop", loopCount)
			// byteIndex = 0
			// count = 1
			// loopCount++
			// if loops > 0 && loopCount >= loops {
			// 	break
			// }
			break
		}

		var rate float64 = -1
		if sendBitRate > 0 {
			/*
			* calculate our instant rate over the entire transmit
			* duration
			 */
			rate = ((float64)(bytesWrittenTotal * 8)) / ((float64)(time.Now().UnixNano()-streamStart) / 1000000000)
		}

		// compare rate
		if rate < sendBitRate {
			max = min(bufferLen, count*int(chunkSize))
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

func (c Connection) SendFile(filepath string, chunkSize int) error {
	addr := c.GetAddr()
	log.Info("Send file", filepath, "to", addr)

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
	conn, err := c.GetSocketConn()
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
		min = int(math.Min(float64(chunkSize), float64(fileSize-int64(bytesRead))))
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
