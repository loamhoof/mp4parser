package mp4parser

import (
	"encoding/binary"
	"io"
)

type containerBox struct {
	offset   int64
	length   uint32
	_type    string
	children []Box
}

func (b *containerBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes := make([]byte, 4)

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	offset := b.offset + 8
	b.children = make([]Box, 0, 1)
	for offset-b.offset < int64(b.length) {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			break
		}

		if _, err := r.Read(bytes); err != nil {
			return err
		}
		l := binary.BigEndian.Uint32(bytes)

		if _, err := r.Read(bytes); err != nil {
			return err
		}
		boxType := string(bytes)

		box := newBox(boxType, offset)
		if err := box.Parse(r); err != nil {
			return err
		}
		b.children = append(b.children, box)

		offset += int64(l)
	}

	return nil
}

func (b *containerBox) Type() string {
	return b._type
}

func (b *containerBox) Offset() int64 {
	return b.offset
}

func (b *containerBox) Length() uint32 {
	return b.length
}

func (b *containerBox) Children() []Box {
	return b.children
}

func (b *containerBox) Data() Pairs {
	return nil
}
