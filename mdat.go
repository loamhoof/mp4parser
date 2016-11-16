package mp4parser

import (
	"io"
)

type mdatBox struct {
	size   uint64
	fields Fields
}

func (b *mdatBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	return nil
}

func (b *mdatBox) Type() string {
	return "mdat"
}

func (b *mdatBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *mdatBox) Size() uint64 {
	return b.size
}

func (b *mdatBox) Children() []Box {
	return []Box{}
}

func (b *mdatBox) Data() Fields {
	return b.fields
}
