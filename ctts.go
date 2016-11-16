package mp4parser

import (
	"encoding/binary"
	"errors"
	"io"
)

type cttsBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *cttsBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes1); err != nil {
		return err
	}
	version := bytes1[0]

	if _, err := r.Seek(3, io.SeekCurrent); err != nil {
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
		entry := make(Pairs, 0, 2)

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entry = append(entry, &Pair{"sample_count", binary.BigEndian.Uint32(bytes4)})

		if version == 0 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			entry = append(entry, &Pair{"sample_offset", binary.BigEndian.Uint32(bytes4)})
		} else if version == 1 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			sampleOffset, read := binary.Varint(bytes4)
			if read <= 0 {
				return errors.New("")
			}
			entry = append(entry, &Pair{"sample_offset", sampleOffset})
		}

		entries[i] = entry
	}

	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *cttsBox) Type() string {
	return "ctts"
}

func (b *cttsBox) Offset() int64 {
	return b.offset
}

func (b *cttsBox) Length() uint32 {
	return b.length
}

func (b *cttsBox) Children() []Box {
	return []Box{}
}

func (b *cttsBox) Data() Pairs {
	return b.data
}
