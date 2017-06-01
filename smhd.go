package mp4parser

import (
	"encoding/binary"
	"io"
)

type SmhdBox struct {
	baseBox
}

func (b *SmhdBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	b.fields = append(b.fields, &Field{"balance", binary.BigEndian.Uint16(b2), offset, 16, 0})

	return nil
}

func (b *SmhdBox) Type() string {
	return "smhd"
}
