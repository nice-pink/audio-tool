package riff

import (
	"encoding/binary"

	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
)

const (
	RIFF_HEADER_SIZE = 12
	WAV_HEADER_SIZE  = 44
	FMT_HEADER_SIZE  = 8
	DATA_HEADER_SIZE = 8
	RIFF_TAG         = "RIFF"
	WAVE_TAG         = "WAVE"
	DATA_TAG         = "data"
	FMT_TAG          = "fmt "
)

type Format struct {
	IsPcm         bool
	Filesize      uint32
	Fmt           uint16
	Channels      uint16
	SampleRate    uint32
	BitPerSample  uint16
	ByteRate      uint32
	BlockAlign    uint16
	ChunkSizeFmt  uint32
	ChunkSizeData uint32
	JunkSize      uint32
}

func IsValid(data []byte) bool {
	if !HasTagId(data) {
		return false
	}

	// IGNORE: check file size
	// dataLen := len(data)
	// fileSize := binary.BigEndian.Uint32(data[4:9])
	// if dataLen != int(fileSize)-8 {
	// 	return false
	// }

	// fmt tag
	fmt := []byte(FMT_TAG)
	if !util.HasTagAtOffset(data, fmt, 12) {
		return false
	}

	// find data block
	offset := GetDataOffset(data)
	log.Debug(offset)
	return offset > 0
}

func HasTagId(data []byte) bool {
	if len(data) < RIFF_HEADER_SIZE {
		return false
	}

	chunkId := []byte(RIFF_TAG)
	if !util.HasTagAtOffset(data, chunkId, 0) {
		return false
	}

	// var fileSize uint32

	riffType := []byte(WAVE_TAG)
	return util.HasTagAtOffset(data, riffType, 8)
}

func Build(size, fileSize uint32) []byte {
	data := make([]byte, size)
	copy(data[0:4], []byte(RIFF_TAG))
	binary.BigEndian.PutUint32(data[4:8], fileSize-8)
	copy(data[8:12], []byte(WAVE_TAG))
	copy(data[12:16], []byte(FMT_TAG))
	if size < 44 {
		log.Error("riff tag too small", size)
		return []byte{}
	}
	copy(data[size-8:size-4], []byte(DATA_TAG))
	binary.BigEndian.PutUint32(data[size-4:size], fileSize-size)
	return data
}

func GetTagSize(data []byte) int64 {
	offset := GetDataOffset(data)
	return int64(offset) + DATA_HEADER_SIZE
}

func GetFileSize(data []byte) uint64 {
	return binary.BigEndian.Uint64(data[4:9])
}

// helper

func GetDataOffset(data []byte) int {
	dataLen := len(data)
	dataTag := []byte(DATA_TAG)
	offset := 36
	for {
		if offset == dataLen {
			// eof
			return -1
		}
		if util.HasTagAtOffset(data, dataTag, offset) {
			// found data
			log.Debug(string(data[offset : offset+4]))
			break
		}
		offset++
	}
	return offset
}
