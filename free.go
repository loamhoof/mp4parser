package mp4parser

import (
	"io"
)

type freeBox struct {
	size   uint64
	_type  string
	fields Fields
}

func (b *freeBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _type, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b._type = _type
	b.fields = fields

	lData := size - uint64(offset-startOffset)
	bData := make([]byte, lData)
	if _, err := r.Read(bData); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"data", string(bData), offset, lData * 8})

	return nil
}

func (b *freeBox) Type() string {
	return b._type
}

func (b *freeBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *freeBox) Size() uint64 {
	return b.size
}

func (b *freeBox) Children() []Box {
	return []Box{}
}

func (b *freeBox) Data() Fields {
	return b.fields
}
