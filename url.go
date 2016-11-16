package mp4parser

import (
	"encoding/binary"
	"io"
)

type urlBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *urlBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if l == 12 {
		b.data = Pairs{{"location", ""}}

		return nil
	}

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}

	bytesLocation := make([]byte, l-12)

	if _, err := r.Read(bytesLocation); err != nil {
		return err
	}
	b.data = Pairs{{"location", string(bytesLocation)}}

	return nil
}

func (b *urlBox) Type() string {
	return "url"
}

func (b *urlBox) Offset() int64 {
	return b.offset
}

func (b *urlBox) Length() uint32 {
	return b.length
}

func (b *urlBox) Children() []Box {
	return []Box{}
}

func (b *urlBox) Data() Pairs {
	return b.data
}
