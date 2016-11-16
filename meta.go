package mp4parser

import (
	"encoding/binary"
	"io"
)

type metaBox struct {
	offset   int64
	length   uint32
	children []Box
}

func (b *metaBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes := make([]byte, 4)

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes)

	b.length = l

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}

	offset := b.offset + 12
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

func (b *metaBox) Type() string {
	return "meta"
}

func (b *metaBox) Offset() int64 {
	return b.offset
}

func (b *metaBox) Length() uint32 {
	return b.length
}

func (b *metaBox) Children() []Box {
	return b.children
}

func (b *metaBox) Data() Pairs {
	return nil
}
