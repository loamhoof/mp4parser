package mp4parser

import (
	"io"
)

type metaBox struct {
	size     uint64
	fields   Fields
	children []Box
}

func (b *metaBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b.children = make([]Box, 0, 1)
	for offset-startOffset < int64(size) {
		size, _, _type, _, err := parseBox(r, offset)
		if err != nil {
			return err
		}

		box := newBox(_type)
		if err := box.Parse(r, offset); err != nil {
			return err
		}
		b.children = append(b.children, box)

		offset += int64(size)
	}

	return nil
}

func (b *metaBox) Type() string {
	return "meta"
}

func (b *metaBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *metaBox) Size() uint64 {
	return b.size
}

func (b *metaBox) Children() []Box {
	return b.children
}

func (b *metaBox) Data() Fields {
	return b.fields
}
