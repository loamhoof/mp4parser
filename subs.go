package mp4parser

import (
	"encoding/binary"
	"io"
)

type subsBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *subsBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes2 := make([]byte, 2)
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
		entry := make(Pairs, 0, 3)

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		entry = append(entry, &Pair{"sample_delta", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes2); err != nil {
			return err
		}
		subsampleCount := binary.BigEndian.Uint16(bytes2)
		entry = append(entry, &Pair{"subsample_count", subsampleCount})

		subsamples := make([]Pairs, subsampleCount)
		for j := 0; uint16(j) < subsampleCount; j++ {
			subsample := make(Pairs, 0, subsampleCount)

			if version == 1 {
				if _, err := r.Read(bytes4); err != nil {
					return err
				}
				subsample = append(subsample, &Pair{"subsample_size", binary.BigEndian.Uint32(bytes4)})
			} else {
				if _, err := r.Read(bytes2); err != nil {
					return err
				}
				subsample = append(subsample, &Pair{"subsample_size", binary.BigEndian.Uint16(bytes2)})
			}

			if _, err := r.Read(bytes1); err != nil {
				return err
			}
			subsample = append(subsample, &Pair{"subsample_priority", bytes1[0]})

			if _, err := r.Read(bytes1); err != nil {
				return err
			}
			subsample = append(subsample, &Pair{"discardable", bytes1[0]})

			if _, err := r.Seek(4, io.SeekCurrent); err != nil {
				return err
			}

			subsamples[j] = subsample
		}

		entry = append(entry, &Pair{"subsamples", subsamples})

		entries[i] = entry
	}
	b.data = append(b.data, &Pair{"entries", entries})

	return nil
}

func (b *subsBox) Type() string {
	return "subs"
}

func (b *subsBox) Offset() int64 {
	return b.offset
}

func (b *subsBox) Length() uint32 {
	return b.length
}

func (b *subsBox) Children() []Box {
	return []Box{}
}

func (b *subsBox) Data() Pairs {
	return b.data
}
