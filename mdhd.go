package mp4parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

type mdhdBox struct {
	size   uint64
	fields Fields
}

func (b *mdhdBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, version, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b2 := make([]byte, 2)
	b4 := make([]byte, 4)
	b8 := make([]byte, 8)

	if version == 1 {
		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"creation_time", binary.BigEndian.Uint64(b8), offset, 64})
		offset += 8

		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"modification_time", binary.BigEndian.Uint64(b8), offset, 64})
		offset += 8

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"timescale", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4

		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"duration", binary.BigEndian.Uint64(b8), offset, 64})
		offset += 8
	} else {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"creation_time", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"modification_time", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"timescale", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"duration", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4
	}

	if _, err := r.Read(b2); err != nil {
		return err
	}
	n := binary.BigEndian.Uint16(b2)

	b.fields = append(b.fields, &Field{"pad", n >> 15, offset, 1}) // TODO
	language := fmt.Sprintf("%c%c%c", n>>10&0x1F+0x60, n>>5&0x1F+0x60, n&0x1F+0x60)
	b.fields = append(b.fields, &Field{"language", language, offset, 15})

	return nil
}

func (b *mdhdBox) Type() string {
	return "mdhd"
}

func (b *mdhdBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *mdhdBox) Size() uint64 {
	return b.size
}

func (b *mdhdBox) Children() []Box {
	return []Box{}
}

func (b *mdhdBox) Data() Fields {
	return b.fields
}
