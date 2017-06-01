package mp4parser

import (
	"encoding/binary"
	"io"
)

type VmhdBox struct {
	baseBox
}

func (b *VmhdBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	b.fields = append(b.fields, &Field{"graphicsmode", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	var opcolor [3]uint16
	for i := 0; i < 3; i++ {
		if _, err := r.Read(b2); err != nil {
			return err
		}
		opcolor[i] = binary.BigEndian.Uint16(b2)
	}
	b.fields = append(b.fields, &Field{"opcolor", opcolor, offset, 48, 0})

	return nil
}

func (b *VmhdBox) Type() string {
	return "vmhd"
}
