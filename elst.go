package mp4parser

import (
	"encoding/binary"
	"io"
)

type elstBox struct {
	baseBox
}

func (b *elstBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, offset, _, version, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b2 := make([]byte, 2)
	b4 := make([]byte, 4)
	b8 := make([]byte, 8)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	entryCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32, 0})
	offset += 4

	entriesOffset := offset
	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := make(Fields, 0, 4)

		if version == 1 {
			if _, err := r.Read(b8); err != nil {
				return err
			}
			entry = append(entry, &Field{"segment_duration", binary.BigEndian.Uint64(b8), offset, 64, 0})
			offset += 8

			if _, err := r.Read(b8); err != nil {
				return err
			}
			entry = append(entry, &Field{"media_time", binary.BigEndian.Uint64(b8), offset, 64, 0})
			offset += 8
		} else {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			entry = append(entry, &Field{"segment_duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4

			if _, err := r.Read(b4); err != nil {
				return err
			}
			entry = append(entry, &Field{"media_time", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4
		}

		if _, err := r.Read(b2); err != nil {
			return err
		}
		entry = append(entry, &Field{"media_rate_integer", binary.BigEndian.Uint16(b2), offset, 16, 0})
		offset += 2

		if _, err := r.Read(b2); err != nil {
			return err
		}
		entry = append(entry, &Field{"media_rate_fraction", binary.BigEndian.Uint16(b2), offset, 16, 0})
		offset += 2

		entries[i] = entry
	}
	b.fields = append(b.fields, &Field{"entries", entries, entriesOffset, uint64(offset-entriesOffset) * 8, 0})

	return nil
}

func (b *elstBox) Type() string {
	return "elst"
}
