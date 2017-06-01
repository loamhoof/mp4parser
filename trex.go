package mp4parser

import (
	"encoding/binary"
	"io"
)

type TrexBox struct {
	baseBox
}

func (b *TrexBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"default_sample_description_index", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"default_sample_duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"is_leading", b4[0] & 0x0C >> 2, offset, 2, 4})
	b.fields = append(b.fields, &Field{"sample_depends_on", b4[0] & 0x03, offset, 2, 6})
	b.fields = append(b.fields, &Field{"sample_is_depended_on", b4[1] & 0xC0 >> 6, offset + 1, 2, 0})
	b.fields = append(b.fields, &Field{"sample_has_redundancy", b4[1] & 0x30 >> 4, offset + 1, 2, 2})
	b.fields = append(b.fields, &Field{"sample_padding_value", b4[1] & 0x0E >> 1, offset + 1, 3, 4})
	b.fields = append(b.fields, &Field{"sample_is_non_sync_sample", b4[1]&0x01 == 1, offset + 1, 1, 7})
	b.fields = append(b.fields, &Field{"sample_degradation_priority", binary.BigEndian.Uint16(b4[2:4]), offset + 2, 16, 0})

	return nil
}

func (b *TrexBox) Type() string {
	return "trex"
}
