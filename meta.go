package mp4parser

import (
	"io"
)

type MetaBox struct {
	baseBox
	children []Box
}

func (b *MetaBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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

		if _, ok := pp[_type]; ok || pp == nil {
			box := newBox(_type)
			if err := box.Parse(r, offset, pp[_type], pc); err != nil {
				return err
			}
			b.children = append(b.children, box)
		}

		offset += int64(size)
	}

	return nil
}

func (b *MetaBox) Type() string {
	return "meta"
}

func (b *MetaBox) Children() []Box {
	return b.children
}
