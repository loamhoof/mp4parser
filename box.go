package mp4parser

import (
	"io"
)

type Box interface {
	Parse(r io.ReadSeeker, offset int64, pp ParsePlan, pc ParseContext) error
	Type() string
	Offset() int64
	Size() uint64
	Children() []Box
	Data() Fields
}

func Background() ParseContext {
	return make(ParseContext)
}

func newBox(_type string) Box {
	switch _type {
	default:
		return &UnknownBox{_type: _type}
	case "ftyp":
		return &FtypBox{}
	case "free", "skip":
		return &FreeBox{}
	case "moov":
		return &MoovBox{}
	case "trak":
		return &TrakBox{}
	case "mvhd":
		return &MvhdBox{}
	case "tkhd":
		return &TkhdBox{}
	case "edts":
		return &EdtsBox{}
	case "elst":
		return &ElstBox{}
	case "mdia":
		return &MdiaBox{}
	case "mdhd":
		return &MdhdBox{}
	case "hdlr":
		return &HdlrBox{}
	case "minf":
		return &MinfBox{}
	case "smhd":
		return &SmhdBox{}
	case "vmhd":
		return &VmhdBox{}
	case "dinf":
		return &DinfBox{}
	case "dref":
		return &DrefBox{}
	case "url":
		return &UrlBox{}
	case "stbl":
		return &StblBox{}
	case "stsd":
		return &StsdBox{}
	case "vide_se":
		return &VisualSampleEntry{}
	case "meta_se":
		return &MetadataSampleEntry{}
	case "avcC":
		return &AvcCBox{}
	case "btrt":
		return &BtrtBox{}
	case "clap":
		return &ClapBox{}
	case "pasp":
		return &PaspBox{}
	case "stts":
		return &SttsBox{}
	case "stss":
		return &StssBox{}
	case "ctts":
		return &CttsBox{}
	case "stsc":
		return &StscBox{}
	case "stsz":
		return &StszBox{}
	case "stco":
		return &StcoBox{}
	case "mvex":
		return &MvexBox{}
	case "mehd":
		return &MehdBox{}
	case "trex":
		return &TrexBox{}
	case "moof":
		return &MoofBox{}
	case "mfhd":
		return &MfhdBox{}
	case "traf":
		return &TrafBox{}
	case "tfhd":
		return &TfhdBox{}
	case "trun":
		return &TrunBox{}
	case "subs":
		return &SubsBox{}
	case "mdat":
		return &MdatBox{}
	case "mfra":
		return &MfraBox{}
	case "tfra":
		return &TfraBox{}
	case "mfro":
		return &MfroBox{}
	case "udta":
		return &UdtaBox{}
	case "meta":
		return &MetaBox{}
	case "ilst":
		return &IlstBox{}
	case "GPMF":
		return &GPMFBox{}
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
