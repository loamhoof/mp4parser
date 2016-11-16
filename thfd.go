package mp4parser

import (
	"encoding/binary"
	"io"
)

type tfhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *tfhdBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes4 := make([]byte, 4)
	bytes8 := make([]byte, 8)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(7, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes1); err != nil {
		return err
	}
	flags := bytes1[0]

	b.data = make(Pairs, 0, 6)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"track_ID", binary.BigEndian.Uint32(bytes4)})

	if flags&0x01 == 0x01 {
		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"base_data_offset", binary.BigEndian.Uint64(bytes8)})
	}

	if flags&0x02 == 0x02 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"sample_description_index", binary.BigEndian.Uint32(bytes4)})
	}

	if flags&0x08 == 0x08 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"default_sample_duration", binary.BigEndian.Uint32(bytes4)})
	}

	if flags&0x10 == 0x10 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"default_sample_size", binary.BigEndian.Uint32(bytes4)})
	}

	if flags&0x20 == 0x20 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"default_sample_flags", binary.BigEndian.Uint32(bytes4)})
	}

	return nil
}

func (b *tfhdBox) Type() string {
	return "tfhd"
}

func (b *tfhdBox) Offset() int64 {
	return b.offset
}

func (b *tfhdBox) Length() uint32 {
	return b.length
}

func (b *tfhdBox) Children() []Box {
	return []Box{}
}

func (b *tfhdBox) Data() Pairs {
	return b.data
}
