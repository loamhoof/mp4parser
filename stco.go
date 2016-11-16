package mp4parser

import (
	"encoding/binary"
	"io"
)

type stcoBox struct {
	size   uint64
	fields Fields
}

func (b *stcoBox) Parse(r io.ReadSeeker, startOffset int64) error {
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

	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		entries[i] = Fields{{"chunk_offset", binary.BigEndian.Uint32(b4), offset, 32}}
		offset += 4
	}
	b.fields = append(b.fields, &Field{"entries", entryCount, offset, 64 * uint64(entryCount)}) // TODO offset

	return nil
}

func (b *stcoBox) Type() string {
	return "stco"
}

func (b *stcoBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *stcoBox) Size() uint64 {
	return b.size
}

func (b *stcoBox) Children() []Box {
	return []Box{}
}

func (b *stcoBox) Data() Fields {
	return b.fields
}
