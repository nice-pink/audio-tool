package riff

import "encoding/binary"

const (
	RIFF_HEADER_SIZE = 12
	WAV_HEADER_SIZE  = 44
	FMT_HEADER_SIZE  = 8
	DATA_HEADER_SIZE = 8
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

	dataLen := len(data)

	// check file size
	fileSize := binary.BigEndian.Uint32(data[4:9])
	if dataLen != int(fileSize)-8 {
		return false
	}

	// fmt tag
	fmt := []byte("fmt ")
	for i, d := range fmt {
		if data[12+i] != d {
			return false
		}
	}

	// find data block
	offset := GetDataOffset(data)

	return offset > 0
}

func HasTagId(data []byte) bool {
	if len(data) < RIFF_HEADER_SIZE {
		return false
	}

	chunkId := []byte("RIFF")
	if !HasTagAtOffset(data, chunkId, 0) {
		return false
	}

	// var chunkSize uint32

	riffType := []byte("WAVE")
	return HasTagAtOffset(data, riffType, 8)
}

func Build(size uint32) []byte {
	return []byte{}
}

func GetTagSize(data []byte) uint32 {
	return 0
}

// helper

func GetDataOffset(data []byte) int {
	dataLen := len(data)
	dataTag := []byte("data")
	offset := 36
	for {
		if offset == dataLen {
			// eof
			return -1
		}
		if HasTagAtOffset(data, dataTag, offset) {
			// found data
			break
		}
		offset++
	}
	return offset
}

func HasTagAtOffset(data, tag []byte, dataOffset int) bool {
	for i, d := range tag {
		if data[dataOffset+i] != d {
			return false
		}
	}
	return false
}
