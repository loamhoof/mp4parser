package mp4parser

import (
	"encoding/binary"
	"io"
)

type trexBox struct {
	size   uint64
	fields Fields
}

func (b *trexBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
	b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"default_sample_description_index", binary.BigEndian.Uint32(b4), offset, 32})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"default_sample_duration", binary.BigEndian.Uint32(b4), offset, 32})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	flags := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"is_leading", uint8(flags >> 24 & 0x0C), offset, 2})
	b.fields = append(b.fields, &Field{"sample_depends_on", uint8(flags >> 24 & 0x03), offset, 2}) // TODO
	b.fields = append(b.fields, &Field{"sample_is_depended_on", uint8(flags >> 20 & 0x0C), offset + 1, 2})
	b.fields = append(b.fields, &Field{"sample_has_redundancy", uint8(flags >> 20 & 0x03), offset + 1, 2})
	b.fields = append(b.fields, &Field{"sample_padding_value", uint8(flags >> 16 & 0x0E), offset + 2, 3})
	b.fields = append(b.fields, &Field{"sample_is_non_sync_sample", flags>>16&0x01 == 1, offset + 2, 1})
	b.fields = append(b.fields, &Field{"sample_degradation_priority", uint16(flags & 0xFFFF), offset + 3, 16})

	return nil
}

func (b *trexBox) Type() string {
	return "trex"
}

func (b *trexBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *trexBox) Size() uint64 {
	return b.size
}

func (b *trexBox) Children() []Box {
	return []Box{}
}

func (b *trexBox) Data() Fields {
	return b.fields
}
