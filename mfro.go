package mp4parser

import (
	"encoding/binary"
	"io"
)

type MfroBox struct {
	baseBox
}

func (b *MfroBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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

func (b *MfroBox) Type() string {
	return "mfro"
}
