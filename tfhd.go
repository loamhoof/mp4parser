package mp4parser

import (
	"encoding/binary"
	"io"
)

type TfhdBox struct {
	baseBox
}

func (b *TfhdBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, _, flags, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)
	b8 := make([]byte, 8)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if flags&0x01 == 0x01 {
		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"base_data_offset", binary.BigEndian.Uint64(b8), offset, 64, 0})
		offset += 8
	}

	if flags&0x02 == 0x02 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"sample_description_index", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	if flags&0x08 == 0x08 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"default_sample_duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	if flags&0x10 == 0x10 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"default_sample_size", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	if flags&0x20 == 0x20 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"default_sample_flags", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	return nil
}

func (b *TfhdBox) Type() string {
	return "tfhd"
}
