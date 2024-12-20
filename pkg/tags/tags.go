package tags

import (
	"github.com/nice-pink/audio-tool/pkg/tags/id3"
	"github.com/nice-pink/audio-tool/pkg/tags/quicktime"
	"github.com/nice-pink/audio-tool/pkg/tags/riff"
)

// interface

type TagBlock interface {
	IsValid(data []byte) bool
	HasTagId(data []byte) bool
	Build(size, fileSize uint32) []byte
	GetTagSize(data []byte) int64
}

// tags

type TagType int

const (
	TagTypeId3 TagType = iota
	TagTypeRiff
	TagTypeQuicktime
	TagTypeUnknown
)

func GetTagType(data []byte) TagType {
	if id3.IsValid(data) {
		return TagTypeId3
	}
	if riff.IsValid(data) {
		return TagTypeRiff
	}
	if quicktime.IsValid(data) {
		return TagTypeQuicktime
	}
	return TagTypeUnknown
}

func GetTagSize(tagType TagType, data []byte) int64 {
	if tagType == TagTypeId3 {
		return id3.GetTagSize(data)
	}
	if tagType == TagTypeRiff {
		return riff.GetTagSize(data)
	}
	if tagType == TagTypeQuicktime {
		return quicktime.GetTagSize(data)
	}
	return -1
}
