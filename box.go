package mp4parser

import (
	"fmt"
	"io"
)

type Box interface {
	Parse(io.ReadSeeker) error
	Type() string
	Offset() int64
	Length() uint32
	Children() []Box
	Data() Pairs
}

type Pair struct {
	Key   string
	Value interface{}
}

func (p *Pair) String() string {
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}

type Pairs []*Pair

func newBox(_type string, offset int64) Box {
	switch _type {
	case "moov", "trak", "mdia", "moof", "mfra", "minf", "dinf", "stbl", "mvex", "traf", "edts", "udta":
		return &containerBox{_type: _type, offset: offset}
	case "mfro":
		return &mfroBox{offset: offset}
	case "free", "skip":
		return &freeBox{offset: offset}
	case "tfra":
		return &tfraBox{offset: offset}
	case "mvhd":
		return &mvhdBox{offset: offset}
	case "ftyp":
		return &ftypBox{offset: offset}
	case "tkhd":
		return &tkhdBox{offset: offset}
	case "mdhd":
		return &mdhdBox{offset: offset}
	case "hdlr":
		return &hdlrBox{offset: offset}
	case "vmhd":
		return &vmhdBox{offset: offset}
	case "dref":
		return &drefBox{offset: offset}
	case "url":
		return &urlBox{offset: offset}
	case "mehd":
		return &mehdBox{offset: offset}
	case "trex":
		return &trexBox{offset: offset}
	case "mfhd":
		return &mfhdBox{offset: offset}
	case "tfhd":
		return &tfhdBox{offset: offset}
	case "trun":
		return &trunBox{offset: offset}
	case "mdat":
		return &mdatBox{offset: offset}
	case "subs":
		return &subsBox{offset: offset}
	case "stts":
		return &sttsBox{offset: offset}
	case "stss":
		return &stssBox{offset: offset}
	case "ctts":
		return &cttsBox{offset: offset}
	case "stsc":
		return &stscBox{offset: offset}
	case "stsz":
		return &stszBox{offset: offset}
	case "stco":
		return &stcoBox{offset: offset}
	case "elst":
		return &elstBox{offset: offset}
	case "smhd":
		return &smhdBox{offset: offset}
	case "meta":
		return &metaBox{offset: offset}
	default:
		return &unknownBox{_type: _type, offset: offset}
	}
}
