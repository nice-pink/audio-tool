package id3

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/nice-pink/audio-tool/pkg/util"
)

const (
	ID3_HEADER_SIZE = 10
	ID3_FOOTER_SIZE = 10
)

func IsValid(data []byte) bool {
	if !HasTagId(data) {
		return false
	}

	offset := 3

	// version: yy < FF
	vData, _ := hex.DecodeString("FFFF")
	for i, d := range vData {
		if data[i+offset] == d {
			return false
		}
	}

	// skip flags: xx
	offset = 7

	// version:= zz < 80
	zzData, _ := hex.DecodeString("80808080")
	for i, d := range zzData {
		if data[i+offset] == d {
			return false
		}
	}

	return true
}

func HasTagId(data []byte) bool {
	// I D 3 yy yy xx zz zz zz zz
	// yy < FF
	// xx -> flags
	// zz < 80

	if len(data) < ID3_HEADER_SIZE {
		return false
	}

	// tag id: 49 44 33
	idData := []byte("ID3")
	for i, d := range idData {
		if data[i] != d {
			return false
		}
	}
	return true
}

func GetTagSize(data []byte) uint32 {
	val := binary.BigEndian.Uint32(data[6:10])

	footerSize := 0
	if HasTagFooter(data) {
		footerSize = ID3_FOOTER_SIZE
	}

	return util.Unsynchsafe(val) + uint32(ID3_HEADER_SIZE) + uint32(footerSize)
}

func HasTagFooter(data []byte) bool {
	mask := uint8(0x10)
	return data[5]&mask == mask
}

func Build(size, fileSize uint32) []byte {
	if size < ID3_HEADER_SIZE {
		return []byte{}
	}

	block := make([]byte, size)

	// tag
	tag := []byte("ID3")
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
