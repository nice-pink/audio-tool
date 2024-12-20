package id3

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/nice-pink/audio-tool/pkg/util"
)

const (
	ID3_HEADER_SIZE = 10
	ID3_FOOTER_SIZE = 10
	ID3_TAG         = "ID3"
)

func IsValid(data []byte) bool {
	if !HasTagId(data) {
		return false
	}

	offset := 3

	// version
	vData, _ := hex.DecodeString("FFFF")
	for i, d := range vData {
		if data[i+offset] == d {
			return false
		}
	}

	// skip flags
	offset = 7

	// version
	zzData, _ := hex.DecodeString("80808080")
	for i, d := range zzData {
		if data[i+offset] == d {
			return false
		}
	}

	return true
}

func HasTagId(data []byte) bool {
	if len(data) < ID3_HEADER_SIZE {
		return false
	}

	// tag id
	idData := []byte(ID3_TAG)
	for i, d := range idData {
		if data[i] != d {
			return false
		}
	}
	return true
}

func GetTagSize(data []byte) int64 {
	val := binary.BigEndian.Uint32(data[6:10])

	footerSize := 0
	if HasTagFooter(data) {
		footerSize = ID3_FOOTER_SIZE
	}

	return int64(util.Unsynchsafe(val)) + int64(ID3_HEADER_SIZE) + int64(footerSize)
}

func HasTagFooter(data []byte) bool {
	if len(data) < ID3_HEADER_SIZE {
		return false
	}
	mask := uint8(0x10)
	return data[5]&mask == mask
}

func Build(size, fileSize uint32) []byte {
	if size < ID3_HEADER_SIZE {
		return []byte{}
	}

	block := make([]byte, size)

	// tag
	tag := []byte(ID3_TAG)
	copy(block[0:3], tag)

	// version
	block[3] = 4
	block[4] = 0

	// flags
	block[5] = 0

	// size
	var synchSize uint32 = util.Synchsafe(size - ID3_HEADER_SIZE)
	sizeData := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeData, synchSize)
	copy(block[6:10], sizeData)

	return block
}
