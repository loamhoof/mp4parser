package mp4parser

import (
	"encoding/binary"
	"io"
)

type trexBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *trexBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}

	b.data = make(Pairs, 0, 10)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"track_ID", binary.BigEndian.Uint32(bytes4)})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"default_sample_description_index", binary.BigEndian.Uint32(bytes4)})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"default_sample_duration", binary.BigEndian.Uint32(bytes4)})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	flags := binary.BigEndian.Uint32(bytes4)
	b.data = append(b.data, &Pair{"is_leading", uint8(flags >> 24 & 0x0C)})
	b.data = append(b.data, &Pair{"sample_depends_on", uint8(flags >> 24 & 0x03)})
	b.data = append(b.data, &Pair{"sample_is_depended_on", uint8(flags >> 20 & 0x0C)})
	b.data = append(b.data, &Pair{"sample_has_redundancy", uint8(flags >> 20 & 0x03)})
	b.data = append(b.data, &Pair{"sample_padding_value", uint8(flags >> 16 & 0x0E)})
	b.data = append(b.data, &Pair{"sample_is_non_sync_sample", flags>>16&0x01 == 1})
	b.data = append(b.data, &Pair{"sample_degradation_priority", uint16(flags & 0xFFFF)})

	return nil
}

func (b *trexBox) Type() string {
	return "trex"
}

func (b *trexBox) Offset() int64 {
	return b.offset
}

func (b *trexBox) Length() uint32 {
	return b.length
}

func (b *trexBox) Children() []Box {
	return []Box{}
}

func (b *trexBox) Data() Pairs {
	return b.data
}
