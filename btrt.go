package mp4parser

import (
	"encoding/binary"
	"io"
)

type BtrtBox struct {
	baseBox
}

func (b *BtrtBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	b.fields = append(b.fields, &Field{"bufferSizeDB", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"maxBitrate", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"avgBitrate", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *BtrtBox) Type() string {
	return "btrt"
}
