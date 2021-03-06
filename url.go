package mp4parser

import (
	"io"
)

type UrlBox struct {
	baseBox
}

func (b *UrlBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	if size == 12 {
		b.fields = append(b.fields, &Field{"location", "", offset, 0, 0})

		return nil
	}

	lLocation := size - uint64(offset-startOffset)
	bLocation := make([]byte, lLocation)

	if _, err := r.Read(bLocation); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"location", string(bLocation), offset, lLocation * 8, 0})

	return nil
}

func (b *UrlBox) Type() string {
	return "url"
}
