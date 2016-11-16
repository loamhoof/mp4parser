package mp4parser

import (
	"encoding/binary"
	"io"
)

type mfroBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *mfroBox) Parse(r io.ReadSeeker) error {
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

	if _, err := r.Read(bytes); err != nil {
		return err
	}

	b.data = Pairs{
		{"size", binary.BigEndian.Uint32(bytes)},
	}

	return nil
}

func (b *mfroBox) Type() string {
	return "mfro"
}

func (b *mfroBox) Offset() int64 {
	return b.offset
}

func (b *mfroBox) Length() uint32 {
	return b.length
}

func (b *mfroBox) Children() []Box {
	return []Box{}
}

func (b *mfroBox) Data() Pairs {
	return b.data
}
