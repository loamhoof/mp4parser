package mp4parser

import (
	"io"
)

type containerBox struct {
	size     uint64
	_type    string
	fields   Fields
	children []Box
}

func (b *containerBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
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

func (b *containerBox) Type() string {
	return b._type
}

func (b *containerBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *containerBox) Size() uint64 {
	return b.size
}

func (b *containerBox) Children() []Box {
	return b.children
}

func (b *containerBox) Data() Fields {
	return b.fields
}
