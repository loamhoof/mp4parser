package mp4parser

import (
	"encoding/binary"
	"io"
)

type smhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *smhdBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes2 := make([]byte, 2)
	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}

	b.data = make(Pairs, 0, 1)

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"balance", binary.BigEndian.Uint16(bytes2)})

	return nil
}

func (b *smhdBox) Type() string {
	return "smhd"
}

func (b *smhdBox) Offset() int64 {
	return b.offset
}

func (b *smhdBox) Length() uint32 {
	return b.length
}

func (b *smhdBox) Children() []Box {
	return []Box{}
}

func (b *smhdBox) Data() Pairs {
	return b.data
}
