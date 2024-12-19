package tags

type TagBlock interface {
	IsValid(data []byte) bool
	HasTagId(data []byte) bool
	Build(size uint32) []byte
	GetTagSize(data []byte) uint32
}
