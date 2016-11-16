package mp4parser

import (
	"encoding/binary"
	"io"
)

type stcoBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *stcoBox) Parse(r io.ReadSeeker) error {
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

	b.data = make(Pairs, 0, 2)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	entryCount := binary.BigEndian.Uint32(bytes4)
	b.data = append(b.data, &Pair{"entry_count", entryCount})

	entries := make([]Pairs, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entries[i] = Pairs{{"chunk_offset", binary.BigEndian.Uint32(bytes4)}}
	}
	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *stcoBox) Type() string {
	return "stco"
}

func (b *stcoBox) Offset() int64 {
	return b.offset
}

func (b *stcoBox) Length() uint32 {
	return b.length
}

func (b *stcoBox) Children() []Box {
	return []Box{}
}

func (b *stcoBox) Data() Pairs {
	return b.data
}
