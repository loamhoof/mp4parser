package mp4parser

import (
	"encoding/binary"
	"io"
)

type drefBox struct {
	offset   int64
	length   uint32
	data     Pairs
	children []Box
}

func (b *drefBox) Parse(r io.ReadSeeker) error {
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

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	entryCount := binary.BigEndian.Uint32(bytes4)

	b.data = Pairs{{"entry_count", entryCount}}

	offset := b.offset + 16
	b.children = make([]Box, 0, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			break
		}

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		l := binary.BigEndian.Uint32(bytes4)

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		boxType := string(bytes4[0:3])

		box := newBox(boxType, offset)
		if err := box.Parse(r); err != nil {
			return err
		}
		b.children = append(b.children, box)

		offset += int64(l)
	}

	return nil
}

func (b *drefBox) Type() string {
	return "dref"
}

func (b *drefBox) Offset() int64 {
	return b.offset
}

func (b *drefBox) Length() uint32 {
	return b.length
}

func (b *drefBox) Children() []Box {
	return b.children
}

func (b *drefBox) Data() Pairs {
	return b.data
}
