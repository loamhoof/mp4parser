package mp4parser

import (
	"encoding/binary"
	"io"
)

type stssBox struct {
	size   uint64
	fields Fields
}

func (b *stssBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32})
	offset += 4

	entriesOffset := offset
	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		entries[i] = Fields{{"sample_number", binary.BigEndian.Uint32(b4), offset, 32}}
		offset += 4
	}
	b.fields = append(b.fields, &Field{"entries", entries, entriesOffset, uint64(offset-entriesOffset) * 8})

	return nil
}

func (b *stssBox) Type() string {
	return "stss"
}

func (b *stssBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *stssBox) Size() uint64 {
	return b.size
}

func (b *stssBox) Children() []Box {
	return []Box{}
}

func (b *stssBox) Data() Fields {
	return b.fields
}
