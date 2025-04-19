package network

import (
	"bufio"
	"errors"
	"math"
	"os"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

type LoopHandler func(loopCount int) error
type SimpleHandler func() error

type StreamBufferStatus struct {
	sendBitRate       float64
	bufferLen         int
	chunkSize         int
	bytesWrittenTotal int64
	streamStart       int64
	loopCount         int
}

func (c *Connection) StreamBuffer(buffer []byte,
	sendBitRate float64,
	chunkSize int,
	endless bool,
	initialFn SimpleHandler,
	loopInitialFn, loopCompletionFn LoopHandler) error {
	// status
	status := &StreamBufferStatus{
		bufferLen:         len(buffer),
		sendBitRate:       sendBitRate,
		chunkSize:         chunkSize,
		bytesWrittenTotal: 0,
		streamStart:       time.Now().UnixNano(),
		loopCount:         0,
	}

	// run loop
	var err error
	if endless {
		// endless
		for {
			// initial function
			if initialFn != nil {
				initialFn()
			}

			err = c.streamBufferLoop(buffer, status, loopInitialFn, loopCompletionFn)
			if err != nil {
				log.Err(err, "stream buffer loop error")
				c.CloseSocket()
				for {
					// retry reconnect socket.
					_, err := c.GetSocketConn()
					if err == nil {
						break
					}
				}
			}
		}
	} else {
		// initial function
		if initialFn != nil {
			initialFn()
		}

		// single run
		err = c.streamBufferLoop(buffer, status, loopInitialFn, loopCompletionFn)
		if err != nil {
			log.Err(err, "stream buffer loop error")
			// socket might be broken -> reinit
			c.socketConn.Close()
			c.socketConn = nil
			for {
				// forever retry to create connection
				_, err = c.GetSocketConn()
				if err == nil {
					break
				}
				time.Sleep(1 * time.Second)
			}

		}
	}

	// final log
	streamStop := time.Now().UnixNano()
	passed := status.streamStart - streamStop
	log.Info("Stopped sending. Bytes:", status.bytesWrittenTotal, ". Seconds:", passed)
	return err
}

func (c *Connection) streamBufferLoop(buffer []byte,
	status *StreamBufferStatus,
	loopInitialFn, loopCompletionFn LoopHandler) error {
	// variables
	var err error
	var byteIndex int = 0
	var bytesWrittenCycle int = 0

	// run loop
	var count int = 1
	for {
		if byteIndex == 0 {
			if loopInitialFn != nil {
				loopInitialFn(status.loopCount)
			}
		} else if byteIndex >= status.bufferLen {
			byteIndex = 0
			count = 1
			status.loopCount++
			if loopCompletionFn != nil {
				loopCompletionFn(status.loopCount)
			}
		}

		var rate float64 = -1
		if status.sendBitRate > 0 {
			// calculate our instant rate over the entire transmit duration
			rate = ((float64)(status.bytesWrittenTotal * 8)) / ((float64)(time.Now().UnixNano()-status.streamStart) / 1000000000)
		}

		// compare rate
		if rate < status.sendBitRate {
			bytesWrittenCycle, err = c.SendBytes(buffer, byteIndex, count, status)
			if err != nil {
				return err
			}
			byteIndex += bytesWrittenCycle

			if c.VerboseLogs {
				log.Info(bytesWrittenCycle, "bytes written")
			}
			count++
		}
	}
}

func (c Connection) SendBytes(buffer []byte, byteIndex, count int, status *StreamBufferStatus) (int, error) {
	max := min(status.bufferLen, count*status.chunkSize)
	dist := max - byteIndex
	// send data
	bytesWrittenCycle, err := c.socketConn.Write(buffer[byteIndex:max])
	if err != nil {
		log.Error(err, "could not send data in loop.")
		return 0, err
	}
	if bytesWrittenCycle <= 0 {
		log.Error("bytes written in cycle", bytesWrittenCycle)
		return 0, errors.New("not all data sent")
	}
	if bytesWrittenCycle != dist {
		log.Error("not all bytes sent. Should", dist, ", did", bytesWrittenCycle)
		return 0, errors.New("not all data sent")
	}
	status.bytesWrittenTotal += int64(bytesWrittenCycle)
	return bytesWrittenCycle, nil
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
		log.Error(err, "file sender can't dial.")
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
			log.Error(err, "could not send data from file", filepath)
			return err
		}
		bytesWrittenTotal += bytesWrittenCycle

		if c.VerboseLogs {
			log.Info(bytesWrittenCycle, "bytes written")
		}
	}

	// final log
	streamStop := time.Now().UnixNano()
	passed := streamStart - streamStop
	log.Info("Stopped sending. Bytes:", bytesWrittenTotal, ". Seconds:", passed)
	return nil
}
