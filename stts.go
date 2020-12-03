package mp4parser

import (
	"encoding/binary"
	"io"
)

type SttsBox struct {
	baseBox
}

func (b *SttsBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	entryCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32, 0})
	offset += 4

	entriesOffset := offset
	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := make(Fields, 0, 2)

		if _, err := r.Read(b4); err != nil {
			return err
		}
		entry = append(entry, &Field{"sample_count", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b4); err != nil {
			return err
		}
		entry = append(entry, &Field{"sample_delta", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		entries[i] = entry
	}
	b.fields = append(b.fields, &Field{"entries", entries, entriesOffset, uint64(offset-entriesOffset) * 8, 0})

	return nil
}

func (b *SttsBox) Type() string {
	return "stts"
}
