package mp4parser

import (
	"encoding/binary"
	"io"
)

type mvhdBox struct {
	size   uint64
	fields Fields
}

func (b *mvhdBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
		b.fields = append(b.fields, &Field{"creation_time", binary.BigEndian.Uint64(b8), offset, 64, 0})
		offset += 8

		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"modification_time", binary.BigEndian.Uint64(b8), offset, 64, 0})
		offset += 8

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"timescale", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b8); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"duration", binary.BigEndian.Uint64(b8), offset, 64, 0})
		offset += 8
	} else {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"creation_time", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"modification_time", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"timescale", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"rate", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b2); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"volume", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Seek(10, io.SeekCurrent); err != nil {
		return err
	}
	offset += 10

	var matrix [9]uint32
	for i := 0; i < 9; i++ {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		matrix[i] = binary.BigEndian.Uint32(b4)
	}
	b.fields = append(b.fields, &Field{"matrix", matrix, offset, 288, 0})
	offset += 36

	if _, err := r.Seek(24, io.SeekCurrent); err != nil {
		return err
	}
	offset += 24

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"next_track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *mvhdBox) Type() string {
	return "mvhd"
}

func (b *mvhdBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *mvhdBox) Size() uint64 {
	return b.size
}

func (b *mvhdBox) Children() []Box {
	return []Box{}
}

func (b *mvhdBox) Data() Fields {
	return b.fields
}
