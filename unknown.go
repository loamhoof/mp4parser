package mp4parser

import (
	"encoding/binary"
	"io"
)

type unknownBox struct {
	_type  string
	offset int64
	length uint32
	data   Pairs
}

func (b *unknownBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	bytesData := make([]byte, l-8)
	if _, err := r.Read(bytesData); err != nil {
		return err
	}
	b.data = Pairs{{"data", string(bytesData)}}

	return nil
}

func (b *unknownBox) Type() string {
	return "[UNKNOWN]" + b._type
}

func (b *unknownBox) Offset() int64 {
	return b.offset
}

func (b *unknownBox) Length() uint32 {
	return b.length
}

func (b *unknownBox) Children() []Box {
	return nil
}

func (b *unknownBox) Data() Pairs {
	return b.data
}
