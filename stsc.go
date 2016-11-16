package mp4parser

import (
	"encoding/binary"
	"io"
)

type stscBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *stscBox) Parse(r io.ReadSeeker) error {
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
		entry := make(Pairs, 0, 3)

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entry = append(entry, &Pair{"first_chunk", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entry = append(entry, &Pair{"samples_per_chunk", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entry = append(entry, &Pair{"sample_description_index", binary.BigEndian.Uint32(bytes4)})

		entries[i] = entry
	}

	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *stscBox) Type() string {
	return "stsc"
}

func (b *stscBox) Offset() int64 {
	return b.offset
}

func (b *stscBox) Length() uint32 {
	return b.length
}

func (b *stscBox) Children() []Box {
	return []Box{}
}

func (b *stscBox) Data() Pairs {
	return b.data
}