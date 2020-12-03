package mp4parser

import (
	"io"
)

type MetadataSampleEntry struct {
	baseBox
	_type    string
	children []Box
	Clap     *ClapBox
	Pasp     *PaspBox
	AvcC     *AvcCBox
	Btrt     *BtrtBox
}

func (b *MetadataSampleEntry) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, _type, fields, err := parseSampleEntry(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b._type = _type
	b.fields = fields

	return nil
}

func (b *MetadataSampleEntry) Type() string {
	return "meta_se:" + b._type
}

func (b *MetadataSampleEntry) Children() []Box {
	return b.children
}
