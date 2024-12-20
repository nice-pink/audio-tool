package quicktime

import (
	"encoding/binary"
	"fmt"

	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
)

type QuicktimeSubtype int

const (
	QuicktimeSubtypeM4A QuicktimeSubtype = iota
	QuicktimeSubtypeInvalid
)

type QuicktimeTag struct {
	TagSize int64
	Subtype QuicktimeSubtype
}

const (
	QUICKTIME_HEADER_SIZE_MIN = 12
	QUICKTIME_TAG             = "ftyp"
	M4A_TAG                   = "M4A "
)

func IsValid(data []byte) bool {
	if !HasTagId(data) {
		return false
	}

	return GetSubtype(data) != QuicktimeSubtypeInvalid
}

func HasTagId(data []byte) bool {

	if len(data) < QUICKTIME_HEADER_SIZE_MIN {
		return false
	}

	// tag id
	idData := []byte(QUICKTIME_TAG)
	return util.HasTagAtOffset(data, idData, 4)
}

func GetTagSize(data []byte) int64 {
	dataSize := len(data)

	// count up all container sizes until sync word for aac is found!
	var size int64 = GetBlockSize(data, 0)
	if len(data) == int(size) {
		return size
	}

	var blockSize int64 = 0
	for {
		if len(data) < int(size) {
			return -1
		}

		blockSize = GetBlockSize(data, uint64(size))
		if blockSize == 0 {
			fmt.Println("No block size")
			break
		}

		if size+blockSize == int64(dataSize) {
			break
		}
		size += blockSize
		// fmt.Println("New size", size)
	}

	return size
}

func Build(size, fileSize uint32) []byte {
	if size < QUICKTIME_HEADER_SIZE_MIN {
		return []byte{}
	}

	block := make([]byte, size)

	sizeData := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeData, QUICKTIME_HEADER_SIZE_MIN)
	copy(block[0:4], sizeData)

	// tag
	tag := []byte(QUICKTIME_TAG)
	copy(block[4:8], tag)

	// sub type
	subtype := []byte(M4A_TAG)
	copy(block[8:12], subtype)

	return block
}

// helper

func GetSubtype(data []byte) QuicktimeSubtype {
	tag := []byte(M4A_TAG)
	if util.HasTagAtOffset(data, tag, 8) {
		return QuicktimeSubtypeM4A
	}
	return QuicktimeSubtypeInvalid
}

func GetQuicktimeTag(data []byte) QuicktimeTag {
	var tag QuicktimeTag
	tag.TagSize = GetTagSize(data)
	tag.Subtype = GetSubtype(data)
	return tag
}

func IntStartsWithAdtsSync(data []byte, offset uint64) bool {
	return util.BytesEqualHexWithMask("FFF0", "FFF6", data[offset:offset+2])
}

func GetBlockSize(data []byte, offset uint64) int64 {
	log.Debug(offset)
	return int64(binary.BigEndian.Uint32(data[offset : offset+4]))
}
