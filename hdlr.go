package mp4parser

import (
	"io"
)

type hdlrBox struct {
	size   uint64
	fields Fields
}

func (b *hdlrBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"handler_type", string(b4), offset, 32})
	offset += 4

	if _, err := r.Seek(12, io.SeekCurrent); err != nil {
		return err
	}
	offset += 12

	lName := size - uint64(offset-startOffset)
	bName := make([]byte, lName)
	if _, err := r.Read(bName); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"name", string(bName), offset, lName * 8})

	return nil
}

func (b *hdlrBox) Type() string {
	return "hdlr"
}

func (b *hdlrBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *hdlrBox) Size() uint64 {
	return b.size
}

func (b *hdlrBox) Children() []Box {
	return []Box{}
}

func (b *hdlrBox) Data() Fields {
	return b.fields
}
