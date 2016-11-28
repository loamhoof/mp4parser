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
	case "ftyp":
		return &ftypBox{}
	case "free", "skip":
		return &freeBox{}
	case "moov":
		return &moovBox{}
	case "trak":
		return &trakBox{}
	case "mvhd":
		return &mvhdBox{}
	case "tkhd":
		return &tkhdBox{}
	case "edts":
		return &edtsBox{}
	case "elst":
		return &elstBox{}
	case "mdia":
		return &mdiaBox{}
	case "mdhd":
		return &mdhdBox{}
	case "hdlr":
		return &hdlrBox{}
	case "minf":
		return &minfBox{}
	case "smhd":
		return &smhdBox{}
	case "vmhd":
		return &vmhdBox{}
	case "dinf":
		return &dinfBox{}
	case "dref":
		return &drefBox{}
	case "url":
		return &urlBox{}
	case "stbl":
		return &stblBox{}
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
	case "mvex":
		return &mvexBox{}
	case "mehd":
		return &mehdBox{}
	case "trex":
		return &trexBox{}
	case "moof":
		return &moofBox{}
	case "mfhd":
		return &mfhdBox{}
	case "traf":
		return &trafBox{}
	case "tfhd":
		return &tfhdBox{}
	case "trun":
		return &trunBox{}
	case "subs":
		return &subsBox{}
	case "mdat":
		return &mdatBox{}
	case "mfra":
		return &mfraBox{}
	case "tfra":
		return &tfraBox{}
	case "mfro":
		return &mfroBox{}
	case "udta":
		return &udtaBox{}
	case "meta":
		return &metaBox{}
	}
}

type baseBox struct {
	size   uint64
	fields Fields
}

func (b *baseBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *baseBox) Size() uint64 {
	return b.size
}

func (b *baseBox) Children() []Box {
	return []Box{}
}

func (b *baseBox) Data() Fields {
	return b.fields
}
