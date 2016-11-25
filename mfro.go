package mp4parser

import (
	"encoding/binary"
	"io"
)

type mfroBox struct {
	size   uint64
	fields Fields
}

func (b *mfroBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
	b.fields = append(b.fields, &Field{"size", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *mfroBox) Type() string {
	return "mfro"
}

func (b *mfroBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *mfroBox) Size() uint64 {
	return b.size
}

func (b *mfroBox) Children() []Box {
	return []Box{}
}

func (b *mfroBox) Data() Fields {
	return b.fields
}
