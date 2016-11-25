package mp4parser

import (
	"encoding/binary"
	"io"
)

type mfhdBox struct {
	size   uint64
	fields Fields
}

func (b *mfhdBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
	b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *mfhdBox) Type() string {
	return "mfhd"
}

func (b *mfhdBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *mfhdBox) Size() uint64 {
	return b.size
}

func (b *mfhdBox) Children() []Box {
	return []Box{}
}

func (b *mfhdBox) Data() Fields {
	return b.fields
}
