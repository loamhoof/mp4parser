package mp4parser

import (
	"encoding/binary"
	"errors"
	"io"
)

type trunBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *trunBox) Parse(r io.ReadSeeker) error {
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

	if _, err := r.Seek(1, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	flags := binary.BigEndian.Uint16(bytes2)

	b.data = make(Pairs, 0, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	sampleCount := binary.BigEndian.Uint32(bytes4)
	b.data = append(b.data, &Pair{"sample_count", sampleCount})

	if flags&0x01 == 0x01 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		dataOffset, read := binary.Varint(bytes4)
		if read <= 0 {
			return errors.New("")
		}
		b.data = append(b.data, &Pair{"data_offset", dataOffset})
	}

	if flags&0x04 == 0x04 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"first_sample_flags", binary.BigEndian.Uint32(bytes4)})
	}

	samples := make([]Pairs, sampleCount)
	for i := 0; uint32(i) < sampleCount; i++ {
		sample := make(Pairs, 0, 4)

		if flags>>8&0x01 == 0x01 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			sample = append(sample, &Pair{"sample_duration", binary.BigEndian.Uint32(bytes4)})
		}

		if flags>>8&0x02 == 0x02 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			sample = append(sample, &Pair{"sample_size", binary.BigEndian.Uint32(bytes4)})
		}

		if flags>>8&0x04 == 0x04 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			sample = append(sample, &Pair{"sample_flags", binary.BigEndian.Uint32(bytes4)})
		}

		if flags>>8&0x08 == 0x08 {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			if version == 0 {
				sample = append(sample, &Pair{"sample_composition_time_offset", binary.BigEndian.Uint32(bytes4)})
			} else {
				sampleCompositionTimeOffset, read := binary.Varint(bytes4)
				if read <= 0 {
					return errors.New("")
				}
				sample = append(sample, &Pair{"sample_composition_time_offset", sampleCompositionTimeOffset})
			}
		}

		samples[i] = sample
	}
	b.data = append(b.data, &Pair{"samples", samples})

	return nil
}

func (b *trunBox) Type() string {
	return "trun"
}

func (b *trunBox) Offset() int64 {
	return b.offset
}

func (b *trunBox) Length() uint32 {
	return b.length
}

func (b *trunBox) Children() []Box {
	return []Box{}
}

func (b *trunBox) Data() Pairs {
	return b.data
}
