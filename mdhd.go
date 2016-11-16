package mp4parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

type mdhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *mdhdBox) Parse(r io.ReadSeeker) error {
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

	b.data = make(Pairs, 0, 6)

	if version == 1 {
		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"creation_time", binary.BigEndian.Uint64(bytes8)})

		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"modification_time", binary.BigEndian.Uint64(bytes8)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"timescale", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"duration", binary.BigEndian.Uint64(bytes8)})
	} else {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"creation_time", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"modification_time", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"timescale", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"duration", binary.BigEndian.Uint32(bytes4)})
	}

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	n := binary.BigEndian.Uint16(bytes2)

	b.data = append(b.data, &Pair{"pad", n >> 15})
	b.data = append(b.data, &Pair{"language", fmt.Sprintf("%c%c%c", n>>10&0x1F+0x60, n>>5&0x1F+0x60, n&0x1F+0x60)})

	return nil
}

func (b *mdhdBox) Type() string {
	return "mdhd"
}

func (b *mdhdBox) Offset() int64 {
	return b.offset
}

func (b *mdhdBox) Length() uint32 {
	return b.length
}

func (b *mdhdBox) Children() []Box {
	return []Box{}
}

func (b *mdhdBox) Data() Pairs {
	return b.data
}
