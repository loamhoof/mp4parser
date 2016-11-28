package mp4parser

import (
	"encoding/binary"
	"io"
)

type mehdBox struct {
	baseBox
}

func (b *mehdBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, version, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)
	b8 := make([]byte, 8)

	if version == 1 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"fragment_duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
	} else {
		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"fragment_duration", binary.BigEndian.Uint64(b8), offset, 64, 0})
	}

	return nil
}

func (b *mehdBox) Type() string {
	return "mehd"
}
