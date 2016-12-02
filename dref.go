package mp4parser

import (
	"encoding/binary"
	"io"
)

type drefBox struct {
	baseBox
	children []Box
	url      []Box
	urn      []Box
}

func (b *drefBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	entryCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32, 0})
	offset += 4

	b.children = make([]Box, 0, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		size, _, _type, _, err := parseBox(r, offset)
		if err != nil {
			return err
		}
		_type = _type[0:3]

		if _, ok := pp[_type]; ok || pp == nil {
			box := newBox(_type)
			if err := box.Parse(r, offset, pp[_type]); err != nil {
				return err
			}
			b.children = append(b.children, box)

			if _type == "url" {
				b.url = append(b.url, box)
			} else if _type == "urn" {
				b.urn = append(b.urn, box)
			}
		}

		offset += int64(size)
	}

	return nil
}

func (b *drefBox) Type() string {
	return "dref"
}

func (b *drefBox) Children() []Box {
	return b.children
}
