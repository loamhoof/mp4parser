package mp4parser

import (
	"io"
)

type unknownBox struct {
	_type  string
	size   uint64
	fields Fields
}

func (b *unknownBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	lData := size - uint64(offset-startOffset)
	bData := make([]byte, lData)
	if _, err := r.Read(bData); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"unknown", string(bData), offset, lData * 8})

	return nil
}

func (b *unknownBox) Type() string {
	return "#" + b._type + "#"
}

func (b *unknownBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *unknownBox) Size() uint64 {
	return b.size
}

func (b *unknownBox) Children() []Box {
	return nil
}

func (b *unknownBox) Data() Fields {
	return b.fields
}
