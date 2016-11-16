package mp4parser

import (
	"encoding/binary"
	"io"
)

type vmhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *vmhdBox) Parse(r io.ReadSeeker) error {
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

	b.data = make(Pairs, 0, 2)

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"graphicsmode", binary.BigEndian.Uint16(bytes2)})

	var opcolor [3]uint16
	for i := 0; i < 3; i++ {
		if _, err := r.Read(bytes2); err != nil {
			return err
		}
		opcolor[i] = binary.BigEndian.Uint16(bytes2)
	}
	b.data = append(b.data, &Pair{"opcolor", opcolor})

	return nil
}

func (b *vmhdBox) Type() string {
	return "vmhd"
}

func (b *vmhdBox) Offset() int64 {
	return b.offset
}

func (b *vmhdBox) Length() uint32 {
	return b.length
}

func (b *vmhdBox) Children() []Box {
	return []Box{}
}

func (b *vmhdBox) Data() Pairs {
	return b.data
}
