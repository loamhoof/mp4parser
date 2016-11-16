package mp4parser

import (
	"encoding/binary"
	"io"
)

type vmhdBox struct {
	size   uint64
	fields Fields
}

func (b *vmhdBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b2 := make([]byte, 2)

	if _, err := r.Read(b2); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"graphicsmode", binary.BigEndian.Uint16(b2), offset, 16})
	offset += 2

	var opcolor [3]uint16
	for i := 0; i < 3; i++ {
		if _, err := r.Read(b2); err != nil {
			return err
		}
		opcolor[i] = binary.BigEndian.Uint16(b2)
	}
	b.fields = append(b.fields, &Field{"opcolor", opcolor, offset, 48})

	return nil
}

func (b *vmhdBox) Type() string {
	return "vmhd"
}

func (b *vmhdBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *vmhdBox) Size() uint64 {
	return b.size
}

func (b *vmhdBox) Children() []Box {
	return []Box{}
}

func (b *vmhdBox) Data() Fields {
	return b.fields
}
