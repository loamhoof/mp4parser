package mp4parser

import (
	"io"
)

type urlBox struct {
	size   uint64
	fields Fields
}

func (b *urlBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	if size == 12 {
		b.fields = append(b.fields, &Field{"location", "", offset, 0})

		return nil
	}

	lLocation := size - uint64(offset-startOffset)
	bLocation := make([]byte, lLocation)

	if _, err := r.Read(bLocation); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"location", string(bLocation), offset, lLocation * 8})

	return nil
}

func (b *urlBox) Type() string {
	return "url"
}

func (b *urlBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *urlBox) Size() uint64 {
	return b.size
}

func (b *urlBox) Children() []Box {
	return []Box{}
}

func (b *urlBox) Data() Fields {
	return b.fields
}
