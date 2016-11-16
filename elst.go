package mp4parser

import (
	"encoding/binary"
	"io"
)

type elstBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *elstBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes2 := make([]byte, 2)
	bytes4 := make([]byte, 4)
	bytes8 := make([]byte, 8)

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
		entry := make(Pairs, 0, 4)

		if version == 1 {
			if _, err := r.Read(bytes8); err != nil {
				return err
			}
			entry = append(entry, &Pair{"segment_duration", binary.BigEndian.Uint64(bytes8)})

			if _, err := r.Read(bytes8); err != nil {
				return err
			}
			entry = append(entry, &Pair{"media_time", binary.BigEndian.Uint64(bytes8)})
		} else {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			entry = append(entry, &Pair{"segment_duration", binary.BigEndian.Uint32(bytes4)})

			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			entry = append(entry, &Pair{"media_time", binary.BigEndian.Uint32(bytes4)})
		}

		if _, err := r.Read(bytes2); err != nil {
			return err
		}
		entry = append(entry, &Pair{"media_rate_integer", binary.BigEndian.Uint16(bytes2)})

		if _, err := r.Read(bytes2); err != nil {
			return err
		}
		entry = append(entry, &Pair{"media_rate_fraction", binary.BigEndian.Uint16(bytes2)})

		entries[i] = entry
	}
	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *elstBox) Type() string {
	return "elst"
}

func (b *elstBox) Offset() int64 {
	return b.offset
}

func (b *elstBox) Length() uint32 {
	return b.length
}

func (b *elstBox) Children() []Box {
	return []Box{}
}

func (b *elstBox) Data() Pairs {
	return b.data
}
