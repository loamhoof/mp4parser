package mp4parser

import (
	"encoding/binary"
	"io"
)

type stszBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *stszBox) Parse(r io.ReadSeeker) error {
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

	b.data = make(Pairs, 0, 3)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	sampleSize := binary.BigEndian.Uint32(bytes4)
	b.data = append(b.data, &Pair{"sample_size", sampleSize})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	sampleCount := binary.BigEndian.Uint32(bytes4)
	b.data = append(b.data, &Pair{"sample_count", sampleCount})

	if sampleSize == 0 {
		samples := make([]Pairs, sampleCount)
		for i := 0; uint32(i) < sampleCount; i++ {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			samples[i] = Pairs{{"entry_size", binary.BigEndian.Uint32(bytes4)}}
		}
		b.data = append(b.data, &Pair{"samples", samples})
	}

	return nil
}

func (b *stszBox) Type() string {
	return "stsz"
}

func (b *stszBox) Offset() int64 {
	return b.offset
}

func (b *stszBox) Length() uint32 {
	return b.length
}

func (b *stszBox) Children() []Box {
	return []Box{}
}

func (b *stszBox) Data() Pairs {
	return b.data
}
