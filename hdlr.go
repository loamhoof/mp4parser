package mp4parser

import (
	"encoding/binary"
	"io"
)

type hdlrBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *hdlrBox) Parse(r io.ReadSeeker) error {
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

	b.data = make(Pairs, 0, 2)

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"handler_type", string(bytes4)})

	if _, err := r.Seek(12, io.SeekCurrent); err != nil {
		return err
	}

	bytesName := make([]byte, l-32)
	if _, err := r.Read(bytesName); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"name", string(bytesName)})

	return nil
}

func (b *hdlrBox) Type() string {
	return "hdlr"
}

func (b *hdlrBox) Offset() int64 {
	return b.offset
}

func (b *hdlrBox) Length() uint32 {
	return b.length
}

func (b *hdlrBox) Children() []Box {
	return []Box{}
}

func (b *hdlrBox) Data() Pairs {
	return b.data
}
