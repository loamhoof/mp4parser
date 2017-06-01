package mp4parser

import (
	"encoding/binary"
	"io"
)

type PaspBox struct {
	baseBox
}

func (b *PaspBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"hSpacing", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"vSpacing", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *PaspBox) Type() string {
	return "pasp"
}
