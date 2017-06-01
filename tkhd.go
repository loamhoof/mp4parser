package mp4parser

import (
	"encoding/binary"
	"io"
)

type TkhdBox struct {
	baseBox
}

func (b *TkhdBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
		b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Seek(4, io.SeekCurrent); err != nil {
			return err
		}
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
		b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Seek(4, io.SeekCurrent); err != nil {
			return err
		}
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}
	offset += 8

	if _, err := r.Read(b2); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"layer", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Read(b2); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"alternate_group", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Read(b2); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"volume", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Seek(2, io.SeekCurrent); err != nil {
		return err
	}
	offset += 2

	var matrix [9]uint32
	for i := 0; i < 9; i++ {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		matrix[i] = binary.BigEndian.Uint32(b4)
	}
	b.fields = append(b.fields, &Field{"matrix", matrix, offset, 288, 0})
	offset += 36

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"width", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"height", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	return nil
}

func (b *TkhdBox) Type() string {
	return "tkhd"
}
