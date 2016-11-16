package mp4parser

import (
	"encoding/binary"
	"io"
)

type mfhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *mfhdBox) Parse(r io.ReadSeeker) error {
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
	b.data = Pairs{{"track_ID", binary.BigEndian.Uint32(bytes4)}}

	return nil
}

func (b *mfhdBox) Type() string {
	return "mfhd"
}

func (b *mfhdBox) Offset() int64 {
	return b.offset
}

func (b *mfhdBox) Length() uint32 {
	return b.length
}

func (b *mfhdBox) Children() []Box {
	return []Box{}
}

func (b *mfhdBox) Data() Pairs {
	return b.data
}
