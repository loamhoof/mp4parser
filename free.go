package mp4parser

import (
	"encoding/binary"
	"io"
)

type freeBox struct {
	offset int64
	length uint32
	_type  string
	data   Pairs
}

func (b *freeBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes := make([]byte, 4)

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes)

	b.length = l

	if _, err := r.Read(bytes); err != nil {
		return err
	}

	b._type = string(bytes)

	bytes = make([]byte, l-8)
	if _, err := r.Read(bytes); err != nil {
		return err
	}

	b.data = Pairs{
		{"data", string(bytes)},
	}

	return nil
}

func (b *freeBox) Type() string {
	return b._type
}

func (b *freeBox) Offset() int64 {
	return b.offset
}

func (b *freeBox) Length() uint32 {
	return b.length
}

func (b *freeBox) Children() []Box {
	return []Box{}
}

func (b *freeBox) Data() Pairs {
	return b.data
}
