package mp4parser

import (
	"encoding/binary"
	"io"
)

type mehdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *mehdBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes4 := make([]byte, 4)
	bytes8 := make([]byte, 8)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes1); err != nil {
		return err
	}
	version := bytes1[0]

	if _, err := r.Seek(3, io.SeekCurrent); err != nil {
		return err
	}

	if version == 1 {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = Pairs{{"fragment_duration", binary.BigEndian.Uint32(bytes4)}}
	} else {
		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = Pairs{{"fragment_duration", binary.BigEndian.Uint64(bytes8)}}
	}

	return nil
}

func (b *mehdBox) Type() string {
	return "mehd"
}

func (b *mehdBox) Offset() int64 {
	return b.offset
}

func (b *mehdBox) Length() uint32 {
	return b.length
}

func (b *mehdBox) Children() []Box {
	return []Box{}
}

func (b *mehdBox) Data() Pairs {
	return b.data
}
