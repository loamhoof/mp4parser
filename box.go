package mp4parser

import (
	"io"
)

type Box interface {
	Parse(r io.ReadSeeker, offset int64) error
	Type() string
	Offset() int64
	Size() uint64
	Children() []Box
	Data() Fields
}

func newBox(_type string) Box {
	switch _type {
	default:
		return &unknownBox{_type: _type}
	case "moov", "trak", "mdia", "moof", "mfra", "minf", "dinf", "stbl", "mvex", "traf", "edts", "udta":
		return &containerBox{_type: _type}
	case "ftyp":
		return &ftypBox{}
	case "free", "skip":
		return &freeBox{}
	case "mvhd":
		return &mvhdBox{}
	case "tkhd":
		return &tkhdBox{}
	case "elst":
		return &elstBox{}
	case "mdhd":
		return &mdhdBox{}
	case "hdlr":
		return &hdlrBox{}
	case "smhd":
		return &smhdBox{}
	case "vmhd":
		return &vmhdBox{}
	case "dref":
		return &drefBox{}
	case "url":
		return &urlBox{}
	case "stts":
		return &sttsBox{}
	case "stss":
		return &stssBox{}
	case "ctts":
		return &cttsBox{}
	case "stsc":
		return &stscBox{}
	case "stsz":
		return &stszBox{}
	case "stco":
		return &stcoBox{}
	case "trex":
		return &trexBox{}
	case "mehd":
		return &mehdBox{}
	case "mfhd":
		return &mfhdBox{}
	case "tfhd":
		return &tfhdBox{}
	case "trun":
		return &trunBox{}
	case "subs":
		return &subsBox{}
	case "mdat":
		return &mdatBox{}
	case "tfra":
		return &tfraBox{}
	case "mfro":
		return &mfroBox{}
	case "meta":
		return &metaBox{}
	}
}
