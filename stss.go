package mp4parser

import (
	"encoding/binary"
	"io"
)

type stssBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *stssBox) Parse(r io.ReadSeeker) error {
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
		entries[i] = Pairs{{"sample_number", binary.BigEndian.Uint32(bytes4)}}
	}
	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *stssBox) Type() string {
	return "stss"
}

func (b *stssBox) Offset() int64 {
	return b.offset
}

func (b *stssBox) Length() uint32 {
	return b.length
}

func (b *stssBox) Children() []Box {
	return []Box{}
}

func (b *stssBox) Data() Pairs {
	return b.data
}
