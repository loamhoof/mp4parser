package mp4parser

import (
	"encoding/binary"
"io"
)

type mdatBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *mdatBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes4 := make([]byte, 4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	return nil
}

func (b *mdatBox) Type() string {
	return "mdat"
}

func (b *mdatBox) Offset() int64 {
	return b.offset
}

func (b *mdatBox) Length() uint32 {
	return b.length
}

func (b *mdatBox) Children() []Box {
	return []Box{}
}

func (b *mdatBox) Data() Pairs {
	return nil
}
